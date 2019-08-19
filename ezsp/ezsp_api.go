// refer to UG100 chapter 3.4

package ezsp

import (
	"fmt"

	"github.com/conthing/utils/common"
)

func ezspApiTrace(format string, v ...interface{}) {
	common.Log.Debugf(format, v...)
}

func littleEndianUint16(b []byte) uint16 {
	return uint16(b[0]) + 256*uint16(b[1])
}

func generalResponseError(response *EzspFrame, cmdID byte) error {
	if response == nil {
		return fmt.Errorf("EZSP cmd 0x%x return nil response", cmdID)
	}
	if response.FrameID == EZSP_INVALID_COMMAND {
		if len(response.Data) != 1 {
			return fmt.Errorf("EZSP cmd 0x%x return invalid command but lenof(0x%x) != 1", cmdID, response.Data)
		}
		return fmt.Errorf("EZSP cmd 0x%x return invalid command ezspStatus(%s)", cmdID, ezspStatusToString(response.Data[0]))
	}
	if response.FrameID != cmdID {
		return fmt.Errorf("EZSP cmd 0x%x response ID(0x%x) not match", cmdID, response.FrameID)
	}
	return nil
}

func generalResponseLengthEqual(response *EzspFrame, cmdID byte, respLen int) error {
	if len(response.Data) != respLen {
		return fmt.Errorf("EZSP cmd 0x%x get invalid response length, expect(%d) get(%d)", cmdID, respLen, len(response.Data))
	}
	return nil
}

func generalResponseLengthNoLessThan(response *EzspFrame, cmdID byte, respLen int) error {
	if len(response.Data) < respLen {
		return fmt.Errorf("EZSP cmd 0x%x get invalid response length, expect(>=%d) get(%d)", cmdID, respLen, len(response.Data))
	}
	return nil
}

func generalResponseEzspStatusSuccess(prefix string, ezspStatus byte) error {
	if ezspStatus != EZSP_SUCCESS {
		return fmt.Errorf("%s get error ezspStatus(%s)", prefix, ezspStatusToString(ezspStatus))
	}
	return nil
}

func EzspVersion(desiredProtocolVersion byte) (protocolVersion byte, stackType byte, stackVersion uint16, err error) {
	response, err := EzspFrameSend(EZSP_VERSION, []byte{desiredProtocolVersion})
	if err == nil {
		err = generalResponseError(response, EZSP_VERSION)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_VERSION, 4)
			if err == nil {
				protocolVersion = response.Data[0]
				stackType = response.Data[1]
				stackVersion = littleEndianUint16(response.Data[2:4])
				if desiredProtocolVersion != protocolVersion {
					err = fmt.Errorf("EzspVersion get unexpected protocolVersion(0x%x) != desired(0x%x)", protocolVersion, desiredProtocolVersion)
					return
				}
				ezspApiTrace("EzspVersion get protocolVersion(0x%x) stackType(0x%x) stackVersion(0x%x)", protocolVersion, stackType, stackVersion)
				//return
			}
		}
	}
	return
}

func EzspGetValue(valueId byte) (ezspStatus byte, value []byte, err error) {
	response, err := EzspFrameSend(EZSP_GET_VALUE, []byte{valueId})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_VALUE)
		if err == nil {
			err = generalResponseLengthNoLessThan(response, EZSP_GET_VALUE, 2)
			if err == nil {
				ezspStatus = response.Data[0]
				valueLength := response.Data[1]
				err = generalResponseLengthEqual(response, EZSP_GET_VALUE, 2+int(valueLength))
				if err == nil {
					value = response.Data[2:]
					ezspApiTrace("EzspGetValue(0x%x) get ezspStatus(%s) value(0x%x)", valueId, ezspStatusToString(ezspStatus), value)
					//return
				}
			}
		}
	}
	return
}

// EzspGetValue API

type EmberVersion struct {
	Build   uint16
	Major   byte
	Minor   byte
	Patch   byte
	Special byte
	VerType byte
}

func EzspGetValue_VERSION_INFO() (emberVersion *EmberVersion, err error) {
	ezspStatus, value, err := EzspGetValue(EZSP_VALUE_VERSION_INFO)
	if err == nil {
		err = generalResponseEzspStatusSuccess("EzspGetValue_VERSION_INFO", ezspStatus)
		if err == nil {
			if len(value) != 7 {
				err = fmt.Errorf("EzspGetValue_VERSION_INFO get invalid value length expect(%d) get(%d)", 7, len(value))
				return
			}
			emberVersion = &EmberVersion{Build: littleEndianUint16(value), Major: value[2], Minor: value[3], Patch: value[4], Special: value[5], VerType: value[6]}
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
			return fmt.Errorf("EZSP cmd 0x%x return invalid command ezspStatus(%s)", EZSP_CALLBACK, ezspStatusToString(response.Data[0]))
		}
		return fmt.Errorf("EZSP_CALLBACK should not have response")
	}
	return
}
