// refer to UG100 chapter 3.4

package ezsp

import (
	"encoding/binary"
	"fmt"

	"github.com/conthing/utils/common"
)

type EmberError struct {
	EmberStatus byte
	OccurAt     string
}
type EzspError struct {
	EzspStatus byte
	OccurAt    string
}

func (e EmberError) Error() string {
	return fmt.Sprintf("%s get error emberStatus(%s)", e.OccurAt, emberStatusToString(e.EmberStatus))
}

func (e EzspError) Error() string {
	return fmt.Sprintf("%s get error ezspStatus(%s)", e.OccurAt, ezspStatusToString(e.EzspStatus))
}

func ezspApiTrace(format string, v ...interface{}) {
	common.Log.Debugf(format, v...)
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
				stackVersion = binary.LittleEndian.Uint16(response.Data[2:])
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

func EzspSetValue(valueId byte, value []byte) (err error) {
	data := []byte{valueId, byte(len(value))}
	data = append(data, value...)
	response, err := EzspFrameSend(EZSP_SET_VALUE, data)
	if err == nil {
		err = generalResponseError(response, EZSP_SET_VALUE)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_VALUE, 1)
			if err == nil {
				ezspStatus := response.Data[0]
				if ezspStatus != EZSP_SUCCESS {
					err = EzspError{ezspStatus, fmt.Sprintf("EzspSetValue(0x%x, 0x%x)", valueId, value)}
					return
				}
				ezspApiTrace("EzspSetValue(0x%x, 0x%x)", valueId, value)
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
				value = binary.LittleEndian.Uint16(response.Data[1:])
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
					err = EzspError{ezspStatus, fmt.Sprintf("EzspSetConfigurationValue(%s, 0x%x)", configIDToName(configId), value)}
					return
				}
				ezspApiTrace("EzspSetConfigurationValue(%s, 0x%x) success", configIDToName(configId), value)
			}
		}
	}
	return
}

func EzspSetPolicy(policyId byte, decisionId byte) (err error) {
	response, err := EzspFrameSend(EZSP_SET_POLICY, []byte{policyId, decisionId})
	if err == nil {
		err = generalResponseError(response, EZSP_SET_POLICY)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_POLICY, 1)
			if err == nil {
				ezspStatus := response.Data[0]
				if ezspStatus != EZSP_SUCCESS {
					err = EzspError{ezspStatus, fmt.Sprintf("EzspSetPolicy(0x%x, 0x%x)", policyId, decisionId)}
					return
				}
				ezspApiTrace("EzspSetPolicy(0x%x, 0x%x) success", policyId, decisionId)
			}
		}
	}
	return
}

func EzspGetEUI64() (eui64 uint64, err error) {
	response, err := EzspFrameSend(EZSP_GET_EUI64, []byte{})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_EUI64)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_GET_EUI64, 8)
			if err == nil {
				eui64 = binary.LittleEndian.Uint64(response.Data)
				ezspApiTrace("EzspGetEUI64 0x%016x", eui64)
			}
		}
	}
	return
}

func EzspSetGpioCurrentConfiguration(portPin byte, cfg byte, out byte) (err error) {
	response, err := EzspFrameSend(EZSP_SET_GPIO_CURRENT_CONFIGURATION, []byte{portPin, cfg, out})
	if err == nil {
		err = generalResponseError(response, EZSP_SET_GPIO_CURRENT_CONFIGURATION)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_GPIO_CURRENT_CONFIGURATION, 1)
			if err == nil {
				ezspStatus := response.Data[0]
				if ezspStatus != EZSP_SUCCESS {
					err = EzspError{ezspStatus, fmt.Sprintf("EzspSetGpioCurrentConfiguration(%d, %d, %d)", portPin, cfg, out)}
					return
				}
				ezspApiTrace("EzspSetGpioCurrentConfiguration(%d, %d, %d) success", portPin, cfg, out)
			}
		}
	}
	return
}

func EzspSetRadioPower(power int8) (err error) {
	response, err := EzspFrameSend(EZSP_SET_RADIO_POWER, []byte{byte(power)})
	if err == nil {
		err = generalResponseError(response, EZSP_SET_RADIO_POWER)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_RADIO_POWER, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, fmt.Sprintf("EzspSetRadioPower(%d)", power)}
					return
				}
				ezspApiTrace("EzspSetRadioPower(%d)", power)
			}
		}
	}
	return
}

func EzspGetMfgToken(tokenId byte) (tokenData []byte, err error) {
	response, err := EzspFrameSend(EZSP_GET_MFG_TOKEN, []byte{tokenId})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_MFG_TOKEN)
		if err == nil {
			err = generalResponseLengthNoLessThan(response, EZSP_GET_MFG_TOKEN, 1)
			if err == nil {
				valueLength := response.Data[0]
				err = generalResponseLengthEqual(response, EZSP_GET_MFG_TOKEN, 1+int(valueLength))
				if err == nil {
					tokenData = response.Data[1:]
					ezspApiTrace("EzspGetMfgToken(0x%x) get tokenData(0x%x)", tokenId, tokenData)
				}
			}
		}
	}
	return
}

func EzspSetMfgToken(tokenId byte, tokenData []byte) (err error) {
	data := []byte{tokenId, byte(len(tokenData))}
	data = append(data, tokenData...)
	response, err := EzspFrameSend(EZSP_SET_MFG_TOKEN, data)
	if err == nil {
		err = generalResponseError(response, EZSP_SET_MFG_TOKEN)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_MFG_TOKEN, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, fmt.Sprintf("EzspSetMfgToken(0x%x, 0x%x)", tokenId, tokenData)}
					return
				}
				ezspApiTrace("EzspSetMfgToken(0x%x, 0x%x)", tokenId, tokenData)
			}
		}
	}
	return
}

func EzspGetToken(tokenId byte) (tokenData []byte, err error) {
	response, err := EzspFrameSend(EZSP_GET_TOKEN, []byte{tokenId})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_TOKEN)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_GET_TOKEN, 9)
			if err == nil {
				emberStatus := response.Data[0]
				tokenData = response.Data[1:]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, fmt.Sprintf("EzspGetToken(0x%x)", tokenId)}
					return
				}
				ezspApiTrace("EzspGetToken(0x%x) get tokenData(0x%x)", tokenId, tokenData)
			}
		}
	}
	return
}

func EzspSetToken(tokenId byte, tokenData []byte) (err error) {
	if len(tokenData) != 8 {
		err = fmt.Errorf("EzspSetToken(0x%x, 0x%x) tokenData lenght != 8", tokenId, tokenData)
		return
	}
	data := []byte{tokenId, byte(len(tokenData))}
	data = append(data, tokenData...)
	response, err := EzspFrameSend(EZSP_SET_TOKEN, data)
	if err == nil {
		err = generalResponseError(response, EZSP_SET_TOKEN)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_TOKEN, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, fmt.Sprintf("EzspSetToken(0x%x, 0x%x)", tokenId, tokenData)}
					return
				}
				ezspApiTrace("EzspSetToken(0x%x, 0x%x)", tokenId, tokenData)
			}
		}
	}
	return
}

type EmberNetworkParameters struct {
	ExtendedPanId uint64
	PanId         uint16
	RadioTxPower  int8
	RadioChannel  byte
	JoinMethod    byte
	NwkManagerId  uint16
	NwkUpdateId   byte
	Channels      uint32
}

func EzspGetNetworkParameters() (nodeType byte, parameters *EmberNetworkParameters, err error) {
	response, err := EzspFrameSend(EZSP_GET_NETWORK_PARAMETERS, []byte{})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_NETWORK_PARAMETERS)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_GET_NETWORK_PARAMETERS, 22)
			if err == nil {
				emberStatus := response.Data[0]
				nodeType = response.Data[1]
				p := EmberNetworkParameters{}
				p.ExtendedPanId = binary.LittleEndian.Uint64(response.Data[2:])
				p.PanId = binary.LittleEndian.Uint16(response.Data[10:])
				p.RadioTxPower = int8(response.Data[12])
				p.RadioChannel = response.Data[13]
				p.JoinMethod = response.Data[14]
				p.NwkManagerId = binary.LittleEndian.Uint16(response.Data[15:])
				p.NwkUpdateId = response.Data[17]
				p.Channels = binary.LittleEndian.Uint32(response.Data[18:])
				parameters = &p

				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspGetNetworkParameters()"}
					return
				}
				ezspApiTrace("EzspGetNetworkParameters() get nodeType(%d) parameters(%+v)", nodeType, *parameters)
			}
		}
	}
	return
}

func EzspNetworkInit() (err error) {
	response, err := EzspFrameSend(EZSP_NETWORK_INIT, []byte{})
	if err == nil {
		err = generalResponseError(response, EZSP_NETWORK_INIT)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_NETWORK_INIT, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspNetworkInit()"}
					return
				}
				ezspApiTrace("EzspNetworkInit()")
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
		emberVersion = &EmberVersion{Build: binary.LittleEndian.Uint16(value), Major: value[2], Minor: value[3], Patch: value[4], Special: value[5], VerType: value[6]}
	}
	return
}

func EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE(size uint16) (err error) {
	return EzspSetValue(EZSP_VALUE_MAXIMUM_INCOMING_TRANSFER_SIZE, []byte{byte(size), byte(size >> 8)})
}
func EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE(size uint16) (err error) {
	return EzspSetValue(EZSP_VALUE_MAXIMUM_OUTGOING_TRANSFER_SIZE, []byte{byte(size), byte(size >> 8)})
}

func EzspSetMfgToken_MFG_PHY_CONFIG(phyConfig uint16) (err error) {
	return EzspSetMfgToken(EZSP_MFG_PHY_CONFIG, []byte{byte(phyConfig), byte(phyConfig >> 8)})
}
func EzspGetMfgToken_MFG_PHY_CONFIG() (phyConfig uint16, err error) {
	value, err := EzspGetMfgToken(EZSP_MFG_PHY_CONFIG)
	if err == nil {
		if len(value) != 2 {
			err = fmt.Errorf("EzspGetMfgToken_MFG_PHY_CONFIG get invalid value length expect(%d) get(%d)", 2, len(value))
			return
		}
		phyConfig = binary.LittleEndian.Uint16(value)
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
