package ezsp

import (
	"fmt"
	"sync"
	"time"

	"github.com/conthing/ezsp/ash"
	"github.com/conthing/utils/common"
)

type EzspFrame struct {
	Sequence byte
	Callback byte // 0-not callback, 1-synchronous callback, 2-asynchronous callback
	FrameID  byte
	Data     []byte
}

var mutex sync.Mutex
var sequence byte
var seqMutex sync.Mutex

// callback 发送到这个ch
var CallbackCh = make(chan *EzspFrame, 128)

// 用sequence做key的数组，存放收到的response时发往的ch
var responseChMap [256]chan *EzspFrame

var EzspFrameTraceOn bool

func ezspFrameTrace(format string, v ...interface{}) {
	if EzspFrameTraceOn {
		common.Log.Debugf(format, v...)
	}
}

func (ezspFrame EzspFrame) String() (s string) {
	s = frameIDToName(ezspFrame.FrameID)
	if ezspFrame.Callback == 2 {
		s += "(async)"
	} else if ezspFrame.Callback == 1 {
		s += "(sync)"
	}
	s += fmt.Sprintf(" 0x%x", ezspFrame.Data)
	return
}

func responseChMapClear(i byte) {
	ch := responseChMap[i]
	if ch != nil {
		select {
		case <-ch:
		default:
		}
		close(ch)
		responseChMap[i] = nil
	}
}

// EzspFrameInitVariables 初始化ezsp frame的一些变量，有些会在ASH的接收处理中用到，
// 应该在 AshReset 成功后再次被调用
func EzspFrameInitVariables() {
	sequence = 0

	// 清空 CallbackCh
	select {
	case <-CallbackCh:
	default:
	}

	for i := range responseChMap {
		responseChMapClear(byte(i))
	}
}

func getSequence() byte {
	seqMutex.Lock()
	seq := sequence
	sequence++
	seqMutex.Unlock()
	return seq
}

func ezspFrameParse(data []byte) (*EzspFrame, error) {
	seq := data[0]
	frmCtrl := data[1]
	frmID := data[2]

	if seq-sequence <= 0x80 { /* seq >= sequence */
		return nil, fmt.Errorf("EZSP frame out of sequence recvseq=%d, sequence=%d", seq, sequence)
	}

	if (frmCtrl & 0xE0) != 0x80 {
		return nil, fmt.Errorf("EZSP not a valid frame ctrl byte 0x%x", frmCtrl)
	}
	if (frmCtrl & 0x1) != 0 {
		return nil, fmt.Errorf("EZSP frame overflow")
	}
	if (frmCtrl & 0x2) != 0 {
		return nil, fmt.Errorf("EZSP frame truncated")
	}
	if (frmCtrl & 0x4) != 0 {
		ezspFrameTrace("EZSP frame callback pending")
	}
	callback := byte((frmCtrl >> 3) & 0x3)
	if callback == 3 {
		return nil, fmt.Errorf("EZSP frame unsupported callback ")
	}

	//检查frmID 和 callback是否匹配
	isCallbackID := isValidCallbackID(frmID)
	if isCallbackID && callback == 0 {
		return nil, fmt.Errorf("EZSP frame callback==%d while ID=%s", callback, frameIDToName(frmID))
	} else if isCallbackID == false && callback != 0 {
		return nil, fmt.Errorf("EZSP frame callback==%d while ID=%s", callback, frameIDToName(frmID))
	}

	return &EzspFrame{Sequence: seq, Callback: callback, FrameID: frmID, Data: data[3:]}, nil
}

// AshRecvImp ASH串口接收处理，运行在串口收发线程中
func AshRecvImp(data []byte) error {
	ash.TransceiverStep = 31
	ezspFrame, err := ezspFrameParse(data)
	if err != nil {
		return fmt.Errorf("EZSP frame parse error: %v", err)
	}
	ash.TransceiverStep = 32
	ezspFrameTrace("EZSP recv < %s", ezspFrame)
	if ezspFrame.Callback == 2 { // async callback 给 CallbackCh
		CallbackCh <- ezspFrame
		return nil
	}
	ash.TransceiverStep = 33
	if ezspFrame.Callback == 1 { // sync callback 也给 CallbackCh，另外发个nil给堵塞的发送函数
		CallbackCh <- ezspFrame
	}
	ash.TransceiverStep = 34
	ch := responseChMap[ezspFrame.Sequence]
	if ch != nil {
		if ezspFrame.Callback == 1 {
			ch <- nil
			return nil
		}
		ash.TransceiverStep = 35
		ch <- ezspFrame
	}
	ash.TransceiverStep = 36
	return nil
}

var SendStep = 0

func EzspFrameSend(frmID byte, data []byte) (*EzspFrame, error) {
	fmt.Printf("***SendStep=%d before\n", SendStep)
	SendStep = 1
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		SendStep = 0
	}()
	seq := getSequence()
	ashFrm := []byte{seq, 0, frmID}
	if data != nil {
		ashFrm = append(ashFrm, data...)
	}

	SendStep = 2
	// 创建接收回复的ch
	responseChMapClear(seq) //如果上一轮sequence发送时超时，有可能没有close
	SendStep = 3
	responseChMap[seq] = make(chan *EzspFrame, 1)
	if responseChMap[seq] == nil {
		return nil, fmt.Errorf("EZSP send %s(seq=%d) failed: make chan failed", frameIDToName(frmID), seq)
	}
	SendStep = 4

	err := ash.AshSend(ashFrm)
	if err != nil {
		responseChMapClear(seq)
		return nil, fmt.Errorf("EZSP send %s(seq=%d) failed: ash send failed: %v", frameIDToName(frmID), seq, err)
	}
	ezspFrameTrace("EZSP send > %s 0x%x", frameIDToName(frmID), data)
	SendStep = 5

	select {
	case response := <-responseChMap[seq]:
		SendStep = 6
		close(responseChMap[seq])
		responseChMap[seq] = nil
		return response, nil
	case <-time.After(time.Millisecond * 15000):
		SendStep = 7
		responseChMapClear(seq)
		SendStep = 8
		return nil, fmt.Errorf("EZSP send %s timeout. ASH env \n%s", frameIDToName(frmID), ash.SprintVariables())
	}
}
