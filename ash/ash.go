package ash

import (
	"fmt"
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

var ashRecvRstackFrame = make(chan byte, 1)
var ashNeedSendProcess = make(chan byte, 16) //AshWrite会被不同的线程调用

var ashRejectCondition = false
var ashImmediatelyAck = false
var ashLastRejectCondition = false

var ashRecvNakFrame = false
var ashRecvErrorFrame []byte

var ashSendTime time.Time //最近一次发送的时间，用来判断超时重发 todo 这个时间使用有问题

//var rxbuffer [8][]byte
//var rxGetPtr byte             /*接收处理的偏移，小于ackNumNext*/
var rxIndexNext = byte(7)     /*下一个接收报文的index，自己报文中的ackNum*/
var rxIndexNextSent = byte(7) /*已经发送出去的ackNum*/

var txbuffer [8][]byte
var txPutPtr byte
var txIndexNext byte       /*下一个发送报文的index，自己报文中的frmNum*/
var txIndexConfirming byte /*正在等待ACK的报文index*/

var AshRecv func([]byte) error

func ashTrace(format string, v ...interface{}) {
	//common.Log.Debugf(format, v...)
}

func inc(index byte) byte {
	return byte((index + 1) & 7)
}

func smallthan(index1 byte, index2 byte) bool {
	return ((index1 - index2) & 7) >= 4
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

func getAckNumForSend() byte { /*发送报文中的ackNum字段，调用此函数后才算ACK过*/
	rxIndexNextSent = rxIndexNext
	return rxIndexNext
}

func needAckFrame() bool {
	/*todo 这里的判断待测试*/
	//return rxIndexNextSent != rxIndexNext
	return smallthan(rxIndexNextSent, rxIndexNext)
}

func sendReady() bool {
	/*已经收到的报文大于发送的acknum*/
	return txIndexNext == txIndexConfirming
}

func getSendBuffer() (ashDataFrame []byte) {
	data := txbuffer[txIndexNext]
	if data != nil {
		control := byte(ASH_CONTROLBYTE_DATA | byte(txIndexNext<<4) | getAckNumForSend())
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
			control := byte(ASH_CONTROLBYTE_DATA | byte(txIndexConfirming<<4) | getAckNumForSend() | ASH_CONTROLBYTE_RETX)
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
	common.Log.Errorf("txIndexNext(%d) < ackNum(%d)", txIndexNext, ackNum)
	return fmt.Errorf("ASH recv ackNum ahead of send frmNum")
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
	if byte(control&0x80) == ASH_CONTROLBYTE_DATA {
		dataFrmPseudoRandom(frame[1:])
		err := ackNumProcess(ackNum)
		if err != nil {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv DAT frame with invalid ackNum: %v < %x", err, frame)
		}

		/*更新frmNumNext*/
		if frmNum == rxIndexNext {
			rxIndexNext = inc(rxIndexNext)
			ashTrace("ASH recv < %x", frame)
			//rxbuffer[frmNum] = frame[1:]
			if AshRecv != nil {
				err = AshRecv(frame[1:])
				if err != nil {
					ashRejectCondition = true
					return err
				}
			}
			ashRejectCondition = false
		} else if smallthan(rxIndexNext, frmNum) {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv discontinuous frame sequence. frmNum=%d, reTx=%v, expect frmNum=%d < %x", frmNum, reTx, rxIndexNext, frame)
		} else {
			if reTx {
				ashImmediatelyAck = true //重发的报文，立刻ACK
				common.Log.Warnf("ASH recv repeative resend frame. frmNum=%d, reTx=%v, expect frmNum=%d < %x", frmNum, reTx, rxIndexNext, frame)
			} else { /*初发的帧比想收的帧序号还要小*/
				ashRejectCondition = true
				return fmt.Errorf("ASH recv frame sequence rollback. frmNum=%d, reTx=%v, expect frmNum=%d < %x", frmNum, reTx, rxIndexNext, frame)
			}
		}
	} else if (byte)(control&0xE0) == ASH_CONTROLBYTE_ACK {
		if len(frame) == 1 {
			err := ackNumProcess(ackNum)
			if err != nil {
				ashRejectCondition = true
				return fmt.Errorf("ASH recv DAT frame with invalid ackNum: %v < %x", err, frame)
			}
			ashTrace("ASH recv ACK frame < %x", frame)
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv ACK frame length error < %x", frame)
		}
	} else if (byte)(control&0xE0) == ASH_CONTROLBYTE_NAK {
		if len(frame) == 1 {
			err := ackNumProcess(ackNum)
			if err != nil {
				ashRejectCondition = true
				return fmt.Errorf("ASH recv DAT frame with invalid ackNum: %v < %x", err, frame)
			}
			ashTrace("ASH recv NAK frame < %x", frame)
			ashRecvNakFrame = true
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv NAK frame length error < %x", frame)
		}
	} else if control == ASH_CONTROLBYTE_RSTACK {
		if len(frame) == 3 {
			ashTrace("ASH recv RSTACK frame < %x", frame)
			if frame[1] != 0x02 {
				ashRejectCondition = true
				return fmt.Errorf("unknown ASH version in RSTACK frame")
			}
			ashRecvRstackFrame <- frame[2]
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv RSTACK frame length error < %x", frame)
		}
	} else if control == ASH_CONTROLBYTE_ERROR {
		if len(frame) == 3 {
			common.Log.Warnf("ASH recv ERROR frame < %x", frame) //todo 测试下ERROR frame的格式
			if frame[1] != 0x02 {
				ashRejectCondition = true
				return fmt.Errorf("unknown ASH version in ERROR frame")
			}
			ashRecvErrorFrame = frame[2:]
		} else {
			ashRejectCondition = true
			return fmt.Errorf("ASH recv ERROR frame length error < %x", frame)
		}
	} else {
		ashRejectCondition = true
		return fmt.Errorf("unknown frame control 0x%x", control)
	}

	return nil
}

func ashAckProcess() {
	if ashRejectCondition == false {
		ashLastRejectCondition = false
	}
	if ashLastRejectCondition == false && ashRejectCondition == true {
		ashLastRejectCondition = true
		err := ashSendNakFrame()
		if err != nil {
			common.Log.Errorf("ASH send NAK frame failed: %v", err)
		}
	} else if needAckFrame() || ashImmediatelyAck {
		err := ashSendAckFrame()
		if err != nil {
			common.Log.Errorf("ASH send ACK frame failed: %v", err)
		} else {
			ashImmediatelyAck = false
		}
	}
}

func ashResendProcess() bool {
	ashDataFrame := getResendBuffer()
	if ashDataFrame != nil {
		ashTrace("ASH resend > %x", ashDataFrame)
		dataFrmPseudoRandom(ashDataFrame[1:])
		err := ashSendFrame(ashDataFrame)
		if err != nil {
			common.Log.Errorf("ASH resend failed: %v", err)
		}
		ashSendTime = time.Now()
		return true
	}
	return false
}

func ashSendProcess() bool {
	if sendReady() {
		ashDataFrame := getSendBuffer()
		if ashDataFrame != nil {
			ashTrace("ASH send > %x", ashDataFrame)
			dataFrmPseudoRandom(ashDataFrame[1:])
			err := ashSendFrame(ashDataFrame)
			if err != nil {
				common.Log.Errorf("ASH send failed: %v", err)
			}
			ashSendTime = time.Now()
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
	frame := []byte{ASH_CONTROLBYTE_ACK | getAckNumForSend()}
	ashTrace("ASH send ACK frame > %x", frame)
	return ashSendFrame(frame)
}
func ashSendNakFrame() error {
	frame := []byte{ASH_CONTROLBYTE_NAK | getAckNumForSend()}
	ashTrace("ASH send NAK frame > %x", frame)
	return ashSendFrame(frame)
}

// ashTransceiver 收发任务
func ashTransceiver(errChan chan error) {
	for {
		resent := false
		select {
		case <-ashNeedSendProcess:
		case <-time.After(time.Millisecond * 50):
			err := AshSerialRecv()
			if err != nil {
				defer func() {
					errChan <- err
				}()
				return
			}

			if ashRecvErrorFrame != nil { //todo 将来改成内部处理
				defer func() {
					errChan <- fmt.Errorf("ASH recv ERROR frame errcode=0x%x", ashRecvErrorFrame[0])
				}()
				return
			}

			if ashRecvNakFrame {
				ashRecvNakFrame = false
				resent = ashResendProcess()
			}
			/*重发和发送ACK的处理，最好在所有收到的报文处理完后进行一次性调用*/
			ashAckProcess()
		}
		if resent == false && time.Now().After(ashSendTime.Add(time.Millisecond*1000)) {
			resent = ashResendProcess()
		}
		_ = ashSendProcess()
	}
}

// AshRead 读取收到的报文
//func AshRead() []byte {
//	ret := rxbuffer[rxGetPtr]
//	if nil != ret {
//		rxbuffer[rxGetPtr] = nil
//		rxGetPtr = inc(rxGetPtr)
//	}
//	return ret
//}

// AshSend 写发送报文缓存
func AshSend(data []byte) error {
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
	common.Log.Debug("ASH Reset")
	err := ashSendCancelByte()
	if err != nil {
		return fmt.Errorf("ASH reset failed: %v", err)
	}

	for i := 0; i < 5; i++ {
		_ = ashSendResetFrame() //不管发送是否成功，没有收到回复就超时重发

		select {
		case rstcode := <-ashRecvRstackFrame:
			common.Log.Debugf("ASH RSTACK 0x%x", rstcode)
			// todo 初始化变量
			rxIndexNext = 0
			return nil
		case <-time.After(time.Millisecond * 3000):
			common.Log.Errorf("ASH Reset miss RSTACK")
		}
	}
	return fmt.Errorf("ASH failed to recv RSTACK after 5 retry")
}

// AshStartTransceiver 开启串口收发线程
func AshStartTransceiver(errChan chan error) {
	go ashTransceiver(errChan)
}
