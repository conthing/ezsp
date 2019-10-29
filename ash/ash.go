package ash

import (
	"fmt"
	"io"
	"time"

	"github.com/conthing/utils/common"
)

const (
	/*ASH协议control byte定义*/
	ASH_CONTROLBYTE_DATA   = byte(0x00)
	ASH_CONTROLBYTE_ACK    = byte(0x80)
	ASH_CONTROLBYTE_NAK    = byte(0xA0)
	ASH_CONTROLBYTE_RST    = byte(0xC0)
	ASH_CONTROLBYTE_RSTACK = byte(0xC1)
	ASH_CONTROLBYTE_ERROR  = byte(0xC2)
	ASH_CONTROLBYTE_RETX   = byte(0x08)
)

var ashResetSuccess = false
var ashRecvRstackFrame = make(chan byte, 1)
var ashNeedSendProcess = make(chan byte, 16) //AshWrite会被不同的线程调用

var ashRejectCondition = false
var ashImmediatelyAck = false
var ashLastRejectCondition = false

var ashRecvNakFrame = false
var ashRecvErrorFrame []byte

//最近一次发送的时间，用来判断超时重发
var ashResendTime *time.Time // todo 这个时间使用有问题 send 和 resend 同时存在时
var ashResendCnt byte

var rxIndexNext byte     /*下一个接收报文的index，自己报文中的ackNum*/
var rxIndexNextSent byte //= byte(7) /*已经发送出去的ackNum*/

var rxbuffer [8][]byte
var rxPutPtr byte

var txbuffer [8][]byte //todo 发送失败怎么清空
var txPutPtr byte
var txIndexNext byte       /*下一个发送报文的index，自己报文中的frmNum*/
var txIndexConfirming byte /*正在等待ACK的报文index*/

var transceiverStep byte

var AshRecv func([]byte) error

var AshTraceOn bool

func ashTrace(format string, v ...interface{}) {
	if AshTraceOn {
		common.Log.Debugf(format, v...)
	}
}

// InitVariables 在AshReset成功后必须调用，恢复原始的状态
func InitVariables() {
	// 清空 ashNeedSendProcess
	select {
	case <-ashNeedSendProcess:
	default:
	}

	ashRejectCondition = false
	ashImmediatelyAck = false
	ashLastRejectCondition = false

	ashRecvNakFrame = false
	ashRecvErrorFrame = nil

	ashResendTime = nil
	ashResendCnt = 0

	rxIndexNext = 0
	rxIndexNextSent = 0 //byte(7) /*已经发送出去的ackNum*/

	for i := range rxbuffer {
		rxbuffer[i] = nil
	}
	rxPutPtr = 0

	for i := range txbuffer {
		txbuffer[i] = nil
	}
	txPutPtr = 0
	txIndexNext = 0       /*下一个发送报文的index，自己报文中的frmNum*/
	txIndexConfirming = 0 /*正在等待ACK的报文index*/

	// 最后 ashResetSuccess 变有效
	ashResetSuccess = true
}

func SprintVariables() (str string) {
	return fmt.Sprintf("transceiverStep=%v \nashResetSuccess=%v \nashResendTime=%v \nashResendCnt=%v \ntxIndexNext=%v \ntxIndexConfirming=%v \ntxPutPtr=%v \ntxbuffer=%v",
		transceiverStep, ashResetSuccess, ashResendTime, ashResendCnt, txIndexNext, txIndexConfirming, txPutPtr, txbuffer)
}

func inc(index byte) byte {
	return byte((index + 1) & 7)
}

func smallthan(index1 byte, index2 byte) bool { // 差值在-1或-2都是smallthan，-3认为是+5了
	return ((index1 - index2) & 7) >= 6
}

func dataFrmPseudoRandom(data []byte) {
	rand := byte(0x42)
	for i := range data {
		data[i] ^= rand
		if (rand & 1) == 0 {
			rand = byte((rand >> 1) & 0x7F)
		} else {
			rand = byte(((rand >> 1) & 0x7F) ^ 0xB8)
		}
	}
}

func getAckNumForData() byte { /*host不能通过DAT进行ACK*/
	//rxIndexNextSent = rxIndexNext
	return rxIndexNextSent
}
func getAckNumForAck() byte { /*发送报文中的ackNum字段，调用此函数后才算ACK过*/ /*host不能通过DAT进行ACK*/
	//rxIndexNextSent = rxIndexNext
	return rxIndexNext //inc(rxIndexNextSent)
}

func needAckFrame() bool {
	return rxIndexNextSent != rxIndexNext
}

func sendReady() bool {
	/*txIndexConfirming使对方报文中最新的acknum，是acked+1，txIndexNext是发送过的+1，txIndexConfirming在追赶txIndexNext*/
	return txIndexNext == txIndexConfirming
}

func getSendBuffer() (ashDataFrame []byte) {
	data := txbuffer[txIndexNext]
	if data != nil {
		control := byte(ASH_CONTROLBYTE_DATA | byte(txIndexNext<<4) | getAckNumForData())
		ashDataFrame = []byte{control}
		ashDataFrame = append(ashDataFrame, data...)
		txIndexNext = inc(txIndexNext)
		return
	}
	return nil
}

func getResendBuffer() (ashDataFrame []byte) {
	if smallthan(txIndexConfirming, txIndexNext) {
		data := txbuffer[txIndexConfirming]
		if data != nil {
			control := byte(ASH_CONTROLBYTE_DATA | byte(txIndexConfirming<<4) | getAckNumForData() | ASH_CONTROLBYTE_RETX)
			ashDataFrame = []byte{control}
			ashDataFrame = append(ashDataFrame, data...)
			return
		}
	}
	return nil
}

func ackNumProcess(ackNum byte) error {
	if !smallthan(txIndexNext, ackNum) { //ackNum > txIndexNext 超前ACK了
		if smallthan(txIndexConfirming, ackNum) {
			for txIndexConfirming != ackNum {
				txbuffer[txIndexConfirming] = nil //已发送成功
				txIndexConfirming = inc(txIndexConfirming)
			}
		}
		return nil
	}
	return fmt.Errorf("ASH recv ackNum(%d) ahead of send frmNum(%d)", ackNum, txIndexNext)
}

// ashRecvFrame 接收报文处理
func ashRecvFrame(frame []byte) error {
	if frame == nil { //表示底层收到非法报文，如crc错误，这里要触发NAK
		ashRejectCondition = true
		return nil
	}

	control := frame[0]
	frmNum := byte((control >> 4) & 7)
	ackNum := byte(control & 7)
	reTx := bool((control & 8) == 8)
	if control == ASH_CONTROLBYTE_RSTACK {
		if len(frame) == 3 {
			ashTrace("ASH recv RSTACK frame < 0x%x", frame)
			if frame[1] != 0x02 {
				ashRejectCondition = true
				return fmt.Errorf("ASH recv unknown version in RSTACK frame")
			}
			ashRecvRstackFrame <- frame[2]
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv RSTACK frame length error < 0x%x", frame)
		}
	} else if control == ASH_CONTROLBYTE_ERROR {
		if len(frame) == 3 {
			common.Log.Warnf("ASH recv ERROR frame < 0x%x", frame) //todo 测试下ERROR frame的格式
			if frame[1] != 0x02 {
				ashRejectCondition = true
				return fmt.Errorf("ASH recv unknown version in ERROR frame")
			}
			ashRecvErrorFrame = frame[2:]
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv ERROR frame length error < 0x%x", frame)
		}
	} else if ashResetSuccess == false { // RSTACK 没收到之前不应该收到其他报文
		return fmt.Errorf("ASH recv other frame before RSTACK < 0x%x", frame)
	} else if byte(control&0x80) == ASH_CONTROLBYTE_DATA {
		dataFrmPseudoRandom(frame[1:])
		err := ackNumProcess(ackNum)
		if err != nil {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv DAT frame with invalid ackNum: %v < 0x%x", err, frame)
		}

		/*更新frmNumNext*/
		if frmNum == rxIndexNext {
			rxIndexNext = inc(rxIndexNext)
			ashTrace("ASH recv < 0x%x", frame)
			rxbuffer[frmNum] = frame[1:]
			//if AshRecv != nil {
			//	err = AshRecv(frame[1:])
			//	if err != nil {
			//		ashRejectCondition = true
			//		return err
			//	}
			//}
			ashRejectCondition = false
			if !smallthan(rxIndexNextSent, rxIndexNext) {
				// rxIndexNext刚增加过，如果rxIndexNextSent-rxIndexNext达到-3，接收报文堆积了3条，需要先处理
				return errAshRecvHandleBusy
			}

		} else if smallthan(rxIndexNext, frmNum) {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv discontinuous frame sequence. frmNum=%d, reTx=%v, expect frmNum=%d < 0x%x", frmNum, reTx, rxIndexNext, frame)
		} else {
			if reTx {
				ashImmediatelyAck = true //重发的报文，立刻ACK
				common.Log.Warnf("ASH recv repeative resend frame. frmNum=%d, reTx=%v, expect frmNum=%d < 0x%x", frmNum, reTx, rxIndexNext, frame)
			} else { /*初发的帧比想收的帧序号还要小*/
				ashRejectCondition = true
				return fmt.Errorf("ASH recv frame sequence rollback. frmNum=%d, reTx=%v, expect frmNum=%d < 0x%x", frmNum, reTx, rxIndexNext, frame)
			}
		}
	} else if (byte)(control&0xE0) == ASH_CONTROLBYTE_ACK {
		if len(frame) == 1 {
			err := ackNumProcess(ackNum)
			if err != nil {
				ashRejectCondition = true
				return fmt.Errorf("ASH recv ACK frame with invalid ackNum: %v < 0x%x", err, frame)
			}
			ashTrace("ASH recv ACK frame < 0x%x", frame)
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv ACK frame length error < 0x%x", frame)
		}
	} else if (byte)(control&0xE0) == ASH_CONTROLBYTE_NAK {
		if len(frame) == 1 {
			err := ackNumProcess(ackNum)
			if err != nil {
				ashRejectCondition = true
				return fmt.Errorf("ASH recv NAK frame with invalid ackNum: %v < 0x%x", err, frame)
			}
			common.Log.Warnf("ASH recv NAK frame < 0x%x", frame)
			ashRecvNakFrame = true
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv NAK frame length error < 0x%x", frame)
		}
	} else {
		ashRejectCondition = true
		return fmt.Errorf("ASH recv unknown frame control 0x%x", control)
	}

	return nil
}

func ashAckProcess() bool {
	if ashRejectCondition == false {
		ashLastRejectCondition = false
	}
	if ashLastRejectCondition == false && ashRejectCondition == true {
		ashLastRejectCondition = true
		err := ashSendNakFrame()
		if err != nil {
			common.Log.Errorf("ASH send NAK frame failed: %v", err)
		}
		return true
	} else if needAckFrame() || ashImmediatelyAck {

		err := ashSendAckFrame()
		if err != nil {
			common.Log.Errorf("ASH send ACK frame failed: %v", err)
		} else {
			for rxIndexNextSent != rxIndexNext {
				if AshRecv != nil && rxbuffer[rxIndexNextSent] != nil {
					err := AshRecv(rxbuffer[rxIndexNextSent])
					if err != nil {
						common.Log.Errorf("ASH recv process failed: %v", err)
						//ashRejectCondition = true
						//return err
					}
				}
				rxIndexNextSent = inc(rxIndexNextSent)
			}
			ashImmediatelyAck = false
		}
		return true
	}
	return false
}

func ashResendProcess() (bool, error) {
	ashDataFrame := getResendBuffer()
	if ashDataFrame != nil {
		if ashResendCnt < 2 {
			ashResendCnt++
			//if ashResendCnt > 1 {
			//	ashTrace("ASH %dth resend ", ashResendCnt)
			//}
			ashTrace("ASH %dth resend > 0x%x", ashResendCnt, ashDataFrame)
			dataFrmPseudoRandom(ashDataFrame[1:])
			err := ashSendFrame(ashDataFrame)
			if err != nil {
				common.Log.Errorf("ASH resend failed: %v", err)
			}
			resendTime := time.Now().Add(time.Second * time.Duration(3+ashResendCnt)) //第一次间隔3秒，第2次4秒...
			ashResendTime = &resendTime
		} else {
			ashResendCnt = 0
			return false, fmt.Errorf("ASH resend exceed max count")
		}
		return true, nil
	}
	return false, nil
}

func ashSendProcess() bool {
	if sendReady() {
		ashDataFrame := getSendBuffer()
		if ashDataFrame != nil {
			ashTrace("ASH send > 0x%x", ashDataFrame)
			dataFrmPseudoRandom(ashDataFrame[1:])
			err := ashSendFrame(ashDataFrame)
			if err != nil {
				common.Log.Errorf("ASH send failed: %v", err)
			}
			resendTime := time.Now().Add(time.Second * 3)
			ashResendTime = &resendTime
			ashResendCnt = 0
			return true
		}
	}
	return false
}

func ashSendResetFrame() error {
	frame := []byte{ASH_CONTROLBYTE_RST}
	ashTrace("ASH send RST frame")
	return ashSendFrame(frame)
}
func ashSendAckFrame() error {
	frame := []byte{ASH_CONTROLBYTE_ACK | getAckNumForAck()}
	ashTrace("ASH send ACK frame > 0x%x", frame)
	return ashSendFrame(frame)
}
func ashSendNakFrame() error {
	frame := []byte{ASH_CONTROLBYTE_NAK | getAckNumForData()}
	ashTrace("ASH send NAK frame > 0x%x", frame)
	return ashSendFrame(frame)
}

// ashTransceiver 收发任务
func ashTransceiver(errChan chan error) {
	AshSerialFlush()
	for {
		resent := false
		acknaksent := false //一次循环发送了ACK就不发DAT了
		transceiverStep = 0
		select {
		case <-ashNeedSendProcess:
		case <-time.After(time.Millisecond * 10):
			transceiverStep = 1
			err := AshSerialRecv()
			transceiverStep = 2
			if err == io.EOF {
				continue
			} else if err != nil {
				errChan <- err
				return
			}
			transceiverStep = 3

			if ashRecvErrorFrame != nil { //todo 将来改成内部处理
				errChan <- fmt.Errorf("ASH recv ERROR frame errcode=0x%x", ashRecvErrorFrame[0])
				return
			}

			if ashResetSuccess { // 没收到RSTACK之前不处理
				if ashRecvNakFrame {
					ashRecvNakFrame = false
					//ashImmediatelyAck = true //resent = ashResendProcess()
				}
				/*重发和发送ACK的处理，最好在所有收到的报文处理完后进行一次性调用*/
				acknaksent = ashAckProcess()
			}
		}
		transceiverStep = 4
		if ashResetSuccess && !acknaksent { // 没收到RSTACK之前不处理
			transceiverStep = 5
			if resent == false && ashResendTime != nil && time.Now().After(*ashResendTime) {
				transceiverStep = 6
				var fatal error
				resent, fatal = ashResendProcess()
				if fatal != nil {
					errChan <- fatal
					return
				}
			}
			transceiverStep = 7
			_ = ashSendProcess()
			transceiverStep = 8
		}
	}
}

// AshSend 写发送报文缓存
func AshSend(data []byte) error {
	if ashResetSuccess != true {
		return fmt.Errorf("ASH RST not finished")
	}
	if txbuffer[txPutPtr] != nil {
		return fmt.Errorf("ASH write overflow")
	}
	txbuffer[txPutPtr] = data //保存发送数据，以备重发
	txPutPtr = inc(txPutPtr)
	ashNeedSendProcess <- 1
	return nil
}

// AshReset 复位NCP
func AshReset() error {
	ashTrace("ASH RST")
	err := ashSendCancelByte()
	if err != nil {
		return fmt.Errorf("ASH RST failed: %v", err)
	}

	ashResetSuccess = false
	// 清空 ashRecvRstackFrame
	select {
	case <-ashRecvRstackFrame:
	default:
	}

	for i := 0; i < 5; i++ {
		_ = ashSendResetFrame() //不管发送是否成功，没有收到回复就超时重发

		select {
		case rstcode := <-ashRecvRstackFrame:
			ashTrace("ASH RSTACK 0x%x", rstcode)
			return nil
		case <-time.After(time.Millisecond * 3000):
			common.Log.Errorf("ASH RST miss RSTACK")
		}
	}
	return fmt.Errorf("ASH failed to recv RSTACK after 5 retry")
}

// AshStartTransceiver 开启串口收发线程，AshReset前就要运行起来
func AshStartTransceiver(recvFunc func([]byte) error, errChan chan error) {
	AshRecv = recvFunc
	go ashTransceiver(errChan)
}
