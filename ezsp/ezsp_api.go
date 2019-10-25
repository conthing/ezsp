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

/** @brief This describes the Initial Security features and requirements that
 *  will be used when forming or joining the network.  */
type EmberInitialSecurityState struct {
	/** This bitmask enumerates which security features should be used, as well
	as the presence of valid data within other elements of the
	::EmberInitialSecurityState data structure.  For more details see the
	::EmberInitialSecurityBitmask. */
	bitmask uint16
	/** This is the pre-configured key that can used by devices when joining the
	 *  network if the Trust Center does not send the initial security data
	 *  in-the-clear.
	 *  For the Trust Center, it will be the global link key and <b>must</b> be set
	 *  regardless of whether joining devices are expected to have a pre-configured
	 *  Link Key.
	 *  This parameter will only be used if the EmberInitialSecurityState::bitmask
	 *  sets the bit indicating ::EMBER_HAVE_PRECONFIGURED_KEY*/
	preconfiguredKey [16]byte
	/** This is the Network Key used when initially forming the network.
	 *  This must be set on the Trust Center.  It is not needed for devices
	 *  joining the network.  This parameter will only be used if the
	 *  EmberInitialSecurityState::bitmask sets the bit indicating
	 *  ::EMBER_HAVE_NETWORK_KEY.  */
	networkKey [16]byte
	/** This is the sequence number associated with the network key.  It must
	 *  be set if the Network Key is set.  It is used to indicate a particular
	 *  of the network key for updating and switching.  This parameter will
	 *  only be used if the ::EMBER_HAVE_NETWORK_KEY is set. Generally it should
	 *  be set to 0 when forming the network; joining devices can ignore
	 *  this value.  */
	networkKeySequenceNumber byte
	/** This is the long address of the trust center on the network that will
	 *  be joined.  It is usually NOT set prior to joining the network and
	 *  instead it is learned during the joining message exchange.  This field
	 *  is only examined if ::EMBER_HAVE_TRUST_CENTER_EUI64 is set in the
	 *  EmberInitialSecurityState::bitmask.  Most devices should clear that
	 *  bit and leave this field alone.  This field must be set when using
	 *  commissioning mode.  It is required to be in little-endian format. */
	preconfiguredTrustCenterEui64 uint64
}

type EmberApsFrame struct {
	/** The application profile ID that describes the format of the message. */
	ProfileId uint16
	/** The cluster ID for this message. */
	ClusterId uint16
	/** The source endpoint. */
	SourceEndpoint byte
	/** The destination endpoint. */
	DestinationEndpoint byte
	/** A bitmask of options from the enumeration above. */
	Options uint16
	/** The group ID for this message, if it is multicast mode. */
	GroupId uint16
	/** The sequence number. */
	Sequence byte
}

type EmberZigbeeNetwork struct {
	Channel       byte
	PanId         uint16
	ExtendedPanId uint64
	AllowingJoin  bool
	StackProfile  byte
	NwkUpdateId   byte
}

func (e EmberError) Error() string {
	return fmt.Sprintf("%s get error emberStatus(%s)", e.OccurAt, emberStatusToString(e.EmberStatus))
}

func (e EzspError) Error() string {
	return fmt.Sprintf("%s get error ezspStatus(%s)", e.OccurAt, ezspStatusToString(e.EzspStatus))
}

var EzspApiTraceOn bool

func ezspApiTrace(format string, v ...interface{}) {
	if EzspApiTraceOn {
		common.Log.Debugf(format, v...)
	}
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
	data := []byte{tokenId}
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

func EzspFormNetwork(para *EmberNetworkParameters) (err error) {
	data := make([]byte, 20)
	binary.LittleEndian.PutUint64(data, para.ExtendedPanId)
	binary.LittleEndian.PutUint16(data[8:], para.PanId)
	data[10] = byte(para.RadioTxPower)
	data[11] = para.RadioChannel
	data[12] = para.JoinMethod
	binary.LittleEndian.PutUint16(data[13:], para.NwkManagerId)
	data[15] = para.NwkUpdateId
	binary.LittleEndian.PutUint32(data[16:], para.Channels)

	response, err := EzspFrameSend(EZSP_FORM_NETWORK, data)
	if err == nil {
		err = generalResponseError(response, EZSP_FORM_NETWORK)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_FORM_NETWORK, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "ezspFormNetwork()"}
					return
				}
				ezspApiTrace("ezspFormNetwork()")
			}
		}
	}
	return
}

func EzspSetInitialSecurityState(state *EmberInitialSecurityState) (err error) {
	data := make([]byte, 43)
	binary.LittleEndian.PutUint16(data, state.bitmask)
	for i, v := range state.preconfiguredKey {
		data[2+i] = v
	}
	for i, v := range state.networkKey {
		data[18+i] = v
	}
	data[34] = state.networkKeySequenceNumber
	binary.LittleEndian.PutUint64(data[35:], state.preconfiguredTrustCenterEui64)

	response, err := EzspFrameSend(EZSP_SET_INITIAL_SECURITY_STATE, data)
	if err == nil {
		err = generalResponseError(response, EZSP_SET_INITIAL_SECURITY_STATE)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_INITIAL_SECURITY_STATE, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspSetInitialSecurityState()"}
					return
				}
				ezspApiTrace("EzspSetInitialSecurityState()")
			}
		}
	}
	return
}

func EzspStartScan(scanType byte, channelMask uint32, duration byte) (err error) {
	response, err := EzspFrameSend(EZSP_START_SCAN, []byte{scanType, byte(channelMask), byte(channelMask >> 8), byte(channelMask >> 16), byte(channelMask >> 24), duration})
	if err == nil {
		err = generalResponseError(response, EZSP_START_SCAN)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_START_SCAN, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspStartScan()"}
					return
				}
				ezspApiTrace("EzspStartScan()")
			}
		}
	}
	return
}

func EzspLookupEui64ByNodeId(nodeId uint16) (eui64 uint64, err error) {
	response, err := EzspFrameSend(EZSP_LOOKUP_EUI64_BY_NODE_ID, []byte{byte(nodeId), byte(nodeId >> 8)})
	if err == nil {
		err = generalResponseError(response, EZSP_LOOKUP_EUI64_BY_NODE_ID)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_LOOKUP_EUI64_BY_NODE_ID, 9)
			if err == nil {
				emberStatus := response.Data[0]
				eui64 = binary.LittleEndian.Uint64(response.Data[1:])
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, fmt.Sprintf("EzspLookupEui64ByNodeId(%04x)", nodeId)}
					return
				}
				ezspApiTrace("EzspLookupEui64ByNodeId(%04x) = 0x%016x", nodeId, eui64)
			}
		}
	}
	return
}

func EzspLookupNodeIdByEui64(eui64 uint64) (nodeId uint16, err error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, eui64)

	response, err := EzspFrameSend(EZSP_LOOKUP_NODE_ID_BY_EUI64, data)
	if err == nil {
		err = generalResponseError(response, EZSP_LOOKUP_NODE_ID_BY_EUI64)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_LOOKUP_NODE_ID_BY_EUI64, 2)
			if err == nil {
				nodeId = binary.LittleEndian.Uint16(response.Data)
				if nodeId == EMBER_NULL_NODE_ID {
					err = fmt.Errorf("EzspLookupNodeIdByEui64(%016x) failed", eui64)
					return
				}
				ezspApiTrace("EzspLookupNodeIdByEui64(%016x) = 0x%04x", eui64, nodeId)
			}
		}
	}
	return
}

func EzspAddressTableEntryIsActive(addressTableIndex byte) (active bool, err error) {
	response, err := EzspFrameSend(EZSP_ADDRESS_TABLE_ENTRY_IS_ACTIVE, []byte{addressTableIndex})
	if err == nil {
		err = generalResponseError(response, EZSP_ADDRESS_TABLE_ENTRY_IS_ACTIVE)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_ADDRESS_TABLE_ENTRY_IS_ACTIVE, 1)
			if err == nil {
				active = (response.Data[0] != 0)
				ezspApiTrace("EzspAddressTableEntryIsActive(%d) = %v", addressTableIndex, active)
			}
		}
	}
	return
}

func EzspGetAddressTableRemoteEui64(addressTableIndex byte) (eui64 uint64, err error) {
	response, err := EzspFrameSend(EZSP_GET_ADDRESS_TABLE_REMOTE_EUI64, []byte{addressTableIndex})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_ADDRESS_TABLE_REMOTE_EUI64)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_GET_ADDRESS_TABLE_REMOTE_EUI64, 8)
			if err == nil {
				eui64 = binary.LittleEndian.Uint64(response.Data)
				ezspApiTrace("EzspGetAddressTableRemoteEui64(%d) = 0x%016x", addressTableIndex, eui64)
			}
		}
	}
	return
}

func EzspGetAddressTableRemoteNodeId(addressTableIndex byte) (nodeID uint16, err error) {
	response, err := EzspFrameSend(EZSP_GET_ADDRESS_TABLE_REMOTE_NODE_ID, []byte{addressTableIndex})
	if err == nil {
		err = generalResponseError(response, EZSP_GET_ADDRESS_TABLE_REMOTE_NODE_ID)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_GET_ADDRESS_TABLE_REMOTE_NODE_ID, 2)
			if err == nil {
				nodeID = binary.LittleEndian.Uint16(response.Data)
				ezspApiTrace("EzspGetAddressTableRemoteNodeId(%d) = 0x%04x", addressTableIndex, nodeID)
			}
		}
	}
	return
}

func EzspSendManyToOneRouteRequest(concentratorType uint16, radius byte) {
	response, err := EzspFrameSend(EZSP_SEND_MANY_TO_ONE_ROUTE_REQUEST, []byte{byte(concentratorType), byte(concentratorType >> 8), radius})
	if err == nil {
		err = generalResponseError(response, EZSP_SEND_MANY_TO_ONE_ROUTE_REQUEST)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SEND_MANY_TO_ONE_ROUTE_REQUEST, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspSendManyToOneRouteRequest()"}
					return
				}
				ezspApiTrace("EzspSendManyToOneRouteRequest()")
			}
		}
	}
	return
}

func EzspPermitJoining(duration byte) (err error) {
	response, err := EzspFrameSend(EZSP_PERMIT_JOINING, []byte{duration})
	if err == nil {
		err = generalResponseError(response, EZSP_PERMIT_JOINING)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_PERMIT_JOINING, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspPermitJoining()"}
					return
				}
				ezspApiTrace("EzspPermitJoining(%d)", duration)
			}
		}
	}
	return
}

func EzspRemoveDevice(destShort uint16, destLong uint64, targetLong uint64) (err error) {
	data := make([]byte, 18)
	binary.LittleEndian.PutUint16(data, destShort)
	binary.LittleEndian.PutUint64(data[2:], destLong)
	binary.LittleEndian.PutUint64(data[10:], targetLong)

	response, err := EzspFrameSend(EZSP_REMOVE_DEVICE, data)
	if err == nil {
		err = generalResponseError(response, EZSP_REMOVE_DEVICE)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_REMOVE_DEVICE, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspRemoveDevice()"}
					return
				}
				ezspApiTrace("EzspRemoveDevice(%x)", targetLong)
			}
		}
	}
	return
}

func EzspSendUnicast(outgoingMessageType byte, indexOrDestination uint16, apsFrame *EmberApsFrame, messageTag byte, message []byte) (sequence byte, err error) {
	data := []byte{
		outgoingMessageType,
		byte(indexOrDestination),
		byte(indexOrDestination >> 8),
		byte(apsFrame.ProfileId),
		byte(apsFrame.ProfileId >> 8),
		byte(apsFrame.ClusterId),
		byte(apsFrame.ClusterId >> 8),
		apsFrame.SourceEndpoint,
		apsFrame.DestinationEndpoint,
		byte(apsFrame.Options),
		byte(apsFrame.Options >> 8),
		byte(apsFrame.GroupId),
		byte(apsFrame.GroupId >> 8),
		apsFrame.Sequence,
		messageTag,
		byte(len(message))}
	data = append(data, message...)
	response, err := EzspFrameSend(EZSP_SEND_UNICAST, data)
	if err == nil {
		err = generalResponseError(response, EZSP_SEND_UNICAST)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SEND_UNICAST, 2)
			if err == nil {
				emberStatus := response.Data[0]
				sequence = response.Data[1]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspSendUnicast()"}
					return
				}
				ezspApiTrace("EzspSendUnicast() return seq=%d", sequence)
			}
		}
	}
	return
}

func EzspLeaveNetwork() (err error) {
	response, err := EzspFrameSend(EZSP_LEAVE_NETWORK, []byte{})
	if err == nil {
		err = generalResponseError(response, EZSP_LEAVE_NETWORK)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_LEAVE_NETWORK, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspLeaveNetwork()"}
					return
				}
				ezspApiTrace("EzspLeaveNetwork()")
			}
		}
	}
	return
}

func EzspSetRadioChannel(channel byte) (err error) {
	response, err := EzspFrameSend(EZSP_SET_RADIO_CHANNEL, []byte{channel})
	if err == nil {
		err = generalResponseError(response, EZSP_SET_RADIO_CHANNEL)
		if err == nil {
			err = generalResponseLengthEqual(response, EZSP_SET_RADIO_CHANNEL, 1)
			if err == nil {
				emberStatus := response.Data[0]
				if emberStatus != EMBER_SUCCESS {
					err = EmberError{emberStatus, "EzspSetRadioChannel()"}
					return
				}
				ezspApiTrace("EzspSetRadioChannel(%d)", channel)
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

func EzspSetValue_EXTENDED_SECURITY_BITMASK(mask uint16) (err error) {
	return EzspSetValue(EZSP_VALUE_EXTENDED_SECURITY_BITMASK, []byte{byte(mask), byte(mask >> 8)})
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
