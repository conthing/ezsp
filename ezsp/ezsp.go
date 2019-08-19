package ezsp

import (
	"github.com/conthing/ezsp/ash"
	"github.com/conthing/utils/common"
)

//TickRunning 定时运行tick
func TickRunning(_ chan error) {

	err := ash.AshReset()
	if err != nil {
		common.Log.Errorf("AshReset failed: %v", err)
	} else {
		EzspFrameInitVariables() // 有些变量在ASH的接收线程里会被使用
		ash.InitVariables()      // 上层的变量初始化完成后，最后调用ASH的变量初始化
		common.Log.Info("AshReset OK")
	}

	protocolVersion, stackType, stackVersion, err := EzspVersion(EZSP_PROTOCOL_VERSION)
	if err != nil {
		common.Log.Errorf("EzspVersion failed: %v", err)
	} else {
		common.Log.Infof("EzspVersion return: %x %x %x", protocolVersion, stackType, stackVersion)
	}

	emberVersion, err := EzspGetValue_VERSION_INFO()
	if err != nil {
		common.Log.Errorf("EzspGetValue_VERSION_INFO failed: %v", err)
	} else {
		common.Log.Infof("EzspGetValue_VERSION_INFO return: %+v", emberVersion)
	}

	err = EzspCallback()
	if err != nil {
		common.Log.Errorf("EzspCallback failed: %v", err)
	} else {
		common.Log.Debugf("EzspCallback return OK")
	}

	for {
		select {
		case cb := <-callbackCh:
			EzspCallbackDispatch(cb)
		}

	}
}

const (
	EZSP_PROTOCOL_VERSION = byte(0x04)
	EZSP_STACK_TYPE_MESH  = byte(0x02)
)
