// refer to UG100 chapter 3.4

package ezsp

import (
	"fmt"

	"github.com/conthing/utils/common"
)

type EzspError struct {
	EzspStatus byte
	OccurAt    string
}

func (e EzspError) Error() string {
	return fmt.Sprintf("%s get error ezspStatus(%s)", e.OccurAt, ezspStatusToString(e.EzspStatus))
}

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

func EzspVersion(desiredProtocolVersion byte) (protocolVersion byte, stackType byte, stackVersion uint16, err error) {
	response, err := EzspFrameSend(EZSP_VERSION, []byte{desiredProtocolVersion})
	if err == nil {
		err = generalResponseError(response, EZSP_VERSION)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_VERSION, 4)
			if err == nil {
				protocolVersion = response.Data[0]
				stackType = response.Data[1]
				stackVersion = littleEndianUint16(response.Data[2:])
				if desiredProtocolVersion != protocolVersion {
					err = fmt.Errorf("EzspVersion get unexpected protocolVersion(0x%x) != desired(0x%x)", protocolVersion, desiredProtocolVersion)
					return
				}
				ezspApiTrace("EzspVersion get protocolVersion(0x%x) stackType(0x%x) stackVersion(0x%x)", protocolVersion, stackType, stackVersion)
			}
		}
	}
	return
}

func EzspGetValue(valueId byte) (value []byte, err error) {
	response, err := EzspFrameSend(EZSP_GET_VALUE, []byte{valueId})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_VALUE)
		if err == nil {
			err = generalResponseLengthNoLessThan(response, EZSP_GET_VALUE, 2)
			if err == nil {
				ezspStatus := response.Data[0]
				valueLength := response.Data[1]
				err = generalResponseLengthEqual(response, EZSP_GET_VALUE, 2+int(valueLength))
				if err == nil {
					if ezspStatus != EZSP_SUCCESS {
						err = EzspError{ezspStatus, fmt.Sprintf("EzspGetValue(0x%x)", valueId)}
						return
					}
					value = response.Data[2:]
					ezspApiTrace("EzspGetValue(0x%x) get value(0x%x)", valueId, value)
				}
			}
		}
	}
	return
}

func EzspGetConfigurationValue(configId byte) (value uint16, err error) {
	response, err := EzspFrameSend(EZSP_GET_CONFIGURATION_VALUE, []byte{configId})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_CONFIGURATION_VALUE)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_GET_CONFIGURATION_VALUE, 3)
			if err == nil {
				ezspStatus := response.Data[0]
				value = littleEndianUint16(response.Data[1:])
				if ezspStatus != EZSP_SUCCESS {
					err = EzspError{ezspStatus, fmt.Sprintf("EzspGetConfigurationValue(%s)", configIDToName(configId))}
					return
				}
				ezspApiTrace("EzspGetConfigurationValue(%s) get 0x%x", configIDToName(configId), value)
			}
		}
	}
	return
}

func EzspSetConfigurationValue(configId byte, value uint16) (err error) {
	response, err := EzspFrameSend(EZSP_SET_CONFIGURATION_VALUE, []byte{configId, byte(value), byte(value >> 8)})
	if err == nil {
		err = generalResponseError(response, EZSP_SET_CONFIGURATION_VALUE)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_CONFIGURATION_VALUE, 1)
			if err == nil {
				ezspStatus := response.Data[0]
				if ezspStatus != EZSP_SUCCESS {
					err = EzspError{ezspStatus, fmt.Sprintf("EzspSetConfigurationValue(0x%x, 0x%x)", configId, value)}
					return
				}
				ezspApiTrace("EzspSetConfigurationValue(0x%x, 0x%x) success", configId, value)
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

func (v EmberVersion) String() (str string) {
	str = fmt.Sprintf("%d.%d.%d.%d build %d",
		v.Major,
		v.Minor,
		v.Patch,
		v.Special,
		v.Build)
	return
}

func EzspGetValue_VERSION_INFO() (emberVersion *EmberVersion, err error) {
	value, err := EzspGetValue(EZSP_VALUE_VERSION_INFO)
	if err == nil {
		if len(value) != 7 {
			err = fmt.Errorf("EzspGetValue_VERSION_INFO get invalid value length expect(%d) get(%d)", 7, len(value))
			return
		}
		emberVersion = &EmberVersion{Build: littleEndianUint16(value), Major: value[2], Minor: value[3], Patch: value[4], Special: value[5], VerType: value[6]}
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
