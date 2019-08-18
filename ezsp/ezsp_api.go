// refer to UG100 chapter 3.4

package ezsp

import (
	"fmt"

	"github.com/conthing/ezsp/ash"
	"github.com/conthing/utils/common"
)

func ezspApiTrace(format string, v ...interface{}) {
	common.Log.Debugf(format, v...)
}

func littleEndianUint16(b []byte) uint16 {
	return uint16(b[0]) + 256*uint16(b[1])
}

func generalResponseError(response *EzspFrame, cmdID byte, respLen int) error {
	if response.FrameID == EZSP_INVALID_COMMAND {
		return fmt.Errorf("EZSP cmd 0x%x return invalid command 0x%x", cmdID, response.Data)
	}
	if response.FrameID != cmdID {
		return fmt.Errorf("EZSP cmd 0x%x response ID(0x%x) not match", cmdID, response.FrameID)
	}
	if len(response.Data) != respLen {
		return fmt.Errorf("EZSP cmd 0x%x get invalid response length: %x", cmdID, response.Data)
	}
	return nil
}

func EzspVersion(desiredProtocolVersion byte) (protocolVersion byte, stackType byte, stackVersion uint16, err error) {
	response, err := EzspFrameSend(EZSP_VERSION, []byte{desiredProtocolVersion})
	if err == nil {
		err = generalResponseError(response, EZSP_VERSION, 4)
		if err == nil {
			protocolVersion = response.Data[0]
			stackType = response.Data[1]
			stackVersion = littleEndianUint16(response.Data[2:4])
			if desiredProtocolVersion != protocolVersion {
				err = fmt.Errorf("EzspVersion get unexpected protocolVersion(0x%x) != desired(0x%x)", protocolVersion, desiredProtocolVersion)
				return
			}
			ezspApiTrace("EzspVersion get protocolVersion(0x%x) stackType(0x%x) stackVersion(0x%x)", protocolVersion, stackType, stackVersion)
			return
		}
	}
	return
}

func EzspCallback() (err error) {
	response, err := EzspFrameSend(EZSP_CALLBACK, []byte{})
	if err == nil {
		if response == nil { //正常应该返回nil，真正的callback从EzspCallbackDispatch处理
			return nil
		}
		if response.FrameID == EZSP_INVALID_COMMAND {
			return fmt.Errorf("EZSP cmd 0x%x return invalid command 0x%x(EzspStatus)", EZSP_CALLBACK, response.Data)
		}
		return fmt.Errorf("EZSP_CALLBACK should not have response")
	}
	return
}

//TickRunning 定时运行tick
func TickRunning(_ chan error) {
	ash.AshRecv = AshRecvImp

	err := ash.AshReset()
	if err != nil {
		common.Log.Errorf("AshReset failed: %v", err)
	} else {
		common.Log.Info("AshReset OK")
	}

	protocolVersion, stackType, stackVersion, err := EzspVersion(4)
	if err != nil {
		common.Log.Errorf("EzspVersion failed: %v", err)
	} else {
		common.Log.Infof("EzspVersion return: %x %x %x", protocolVersion, stackType, stackVersion)
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
