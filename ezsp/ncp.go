package ezsp

import (
	"fmt"

	"encoding/binary"

	"github.com/conthing/utils/common"
)

var NcpTraceOn bool

func ncpTrace(format string, v ...interface{}) {
	if NcpTraceOn {
		common.Log.Debugf(format, v...)
	}
}

type StNcpCallbacks struct {
	NcpMessageSentHandler         func(outgoingMessageType byte, indexOrDestination uint16, apsFrame *EmberApsFrame, messageTag byte, emberStatus byte, message []byte)
	NcpIncomingSenderEui64Handler func(senderEui64 uint64)
	NcpIncomingMessageHandler     func(incomingMessageType byte, apsFrame *EmberApsFrame, lastHopLqi byte, lastHopRssi int8, sender uint16, bindingIndex byte, addressIndex byte, message []byte)
	NcpTrustCenterJoinHandler     func(newNodeId uint16, newNodeEui64 uint64, deviceUpdateStatus byte, joinDecision byte, parentOfNewNode uint16)
}

// StModuleInfo
type StModuleInfo struct {
	ModuleType      string `json:"moduletype"`
	ProtocolVersion byte   `json:"protocolversion"`
	StackType       byte   `json:"stacktype"`
	StackVersion    string `json:"stackversion"`
}

// StMeshInfo
type StMeshInfo struct {
	ExPANID uint64 `json:"expanid"`
	PANID   uint16 `json:"panid"`
	Channel byte   `json:"channel"`
}

var ModuleInfo = StModuleInfo{ModuleType: "EM357"}
var MeshInfo StMeshInfo
var MeshStatusUp bool
var NcpCallbacks StNcpCallbacks

func NcpGetVersion() (err error) {
	var stackVersion uint16
	ModuleInfo.ProtocolVersion, ModuleInfo.StackType, stackVersion, err = EzspVersion(EZSP_PROTOCOL_VERSION)
	if err != nil {
		return fmt.Errorf("EzspVersion failed: %v", err)
	}

	emberVersion, err := EzspGetValue_VERSION_INFO()
	if err != nil {
		common.Log.Errorf("EzspGetValue_VERSION_INFO failed: %v", err)
		ModuleInfo.StackVersion = fmt.Sprintf("%d.%d.%d.%d", (stackVersion>>12)&0xF, (stackVersion>>8)&0xF, (stackVersion>>4)&0xF, stackVersion&0xF)
	} else {
		ModuleInfo.StackVersion = emberVersion.String()
	}

	//common.Log.Infof("%v", stackVersion)

	ncpTrace("NcpGetVersion: protocolVersion(%d) stackType(%d) stackVersion(%s)", ModuleInfo.ProtocolVersion, ModuleInfo.StackType, ModuleInfo.StackVersion)
	return nil
}

func NcpPrintAllConfigurations() {
	for id := 0; id < 256; id++ {
		name, ok := configIDNameMap[byte(id)]
		if ok {
			value, err := EzspGetConfigurationValue(byte(id))
			if err != nil {
				common.Log.Errorf("%s read failed: %v", name, err)
			}
			ncpTrace("%s = %d", name, value)
		}
	}
}

type EzspConfig struct {
	configID byte
	value    uint16
}

var ncpConfigurations = [...]EzspConfig{
	{EZSP_CONFIG_STACK_PROFILE, uint16(2)},
	{EZSP_CONFIG_SUPPORTED_NETWORKS, uint16(1)},
	{EZSP_CONFIG_ADDRESS_TABLE_SIZE, uint16(64)},
	{EZSP_CONFIG_INDIRECT_TRANSMISSION_TIMEOUT, uint16(7680)},
	{EZSP_CONFIG_PACKET_BUFFER_COUNT, uint16(75)},
	{EZSP_CONFIG_MULTICAST_TABLE_SIZE, uint16(1)},
	{EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT, uint16(255)},
	{EZSP_CONFIG_MOBILE_NODE_POLL_TIMEOUT, uint16(255)},

	//{EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE, uint16(2)},
}

func NcpConfig() (err error) {
	for _, cfg := range ncpConfigurations {
		err = EzspSetConfigurationValue(cfg.configID, cfg.value)
		name := configIDToName(cfg.configID)
		if err != nil {
			return fmt.Errorf("%s write %d failed: %v", name, cfg.value, err)
		}
		value, err := EzspGetConfigurationValue(cfg.configID)
		if err != nil {
			return fmt.Errorf("%s read failed: %v", name, err)
		}
		if value != cfg.value {
			return fmt.Errorf("%s read back %d != %d", name, value, cfg.value)
		}
		ncpTrace("Set %s = %d", name, cfg.value)
	}

	err = EzspSetPolicy(EZSP_MESSAGE_CONTENTS_IN_CALLBACK_POLICY, EZSP_MESSAGE_TAG_AND_CONTENTS_IN_CALLBACK)
	if err != nil {
		return fmt.Errorf("EzspSetPolicy failed: %v", err)
	}
	ncpTrace("EzspSetPolicy EZSP_MESSAGE_TAG_AND_CONTENTS_IN_CALLBACK")

	err = EzspSetPolicy(EZSP_UNICAST_REPLIES_POLICY, EZSP_HOST_WILL_NOT_SUPPLY_REPLY)
	if err != nil {
		return fmt.Errorf("EzspSetPolicy failed: %v", err)
	}
	ncpTrace("EzspSetPolicy EZSP_HOST_WILL_NOT_SUPPLY_REPLY")

	err = EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE(84)
	if err != nil {
		return fmt.Errorf("EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE failed: %v", err)
	}
	ncpTrace("EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE = 84")

	err = EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE(84)
	if err != nil {
		return fmt.Errorf("EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE failed: %v", err)
	}
	ncpTrace("EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE = 84")

	err = ncpSetRadio()
	if err != nil {
		return fmt.Errorf("ncpSetRadio failed: %v", err)
	}
	ncpTrace("ncpSetRadio OK")

	return
}

func ncpSetRadio() (err error) {
	err = EzspSetGpioCurrentConfiguration(PORTA_PIN7, 1, 0)
	if err != nil {
		return fmt.Errorf("EzspSetGpioCurrentConfiguration(PORTA_PIN7,1,0) failed: %v", err)
	}
	err = EzspSetGpioCurrentConfiguration(PORTA_PIN3, 1, 1)
	if err != nil {
		return fmt.Errorf("EzspSetGpioCurrentConfiguration(PORTA_PIN3,1,1) failed: %v", err)
	}
	err = EzspSetGpioCurrentConfiguration(PORTA_PIN6, 1, 1)
	if err != nil {
		return fmt.Errorf("EzspSetGpioCurrentConfiguration(PORTA_PIN6,1,1) failed: %v", err)
	}
	err = EzspSetGpioCurrentConfiguration(PORTC_PIN5, 9, 0)
	if err != nil {
		return fmt.Errorf("EzspSetGpioCurrentConfiguration(PORTC_PIN5,9,0) failed: %v", err)
	}

	err = EzspSetRadioPower(3)
	if err != nil {
		return fmt.Errorf("ezspSetRadioPower(3) failed: %v", err)
	}

	phyConfig, err := EzspGetMfgToken_MFG_PHY_CONFIG()
	if err != nil {
		return fmt.Errorf("EzspGetMfgToken_MFG_PHY_CONFIG() failed: %v", err)
	}

	if phyConfig != 0xfffd {
		err = EzspSetMfgToken_MFG_PHY_CONFIG(0xfffd)
		if err != nil {
			return fmt.Errorf("EzspSetMfgToken_MFG_PHY_CONFIG(0xfffd) failed: %v", err)
		}
	}

	//只有第一次写入不抱错，以后写都会报次错误
	return nil
}

func NcpGetAndIncRebootCnt() (rebootCnt uint16, err error) {
	//tokenId=0的8个字节定义成NCP使用，低2字节为rebootCnt
	tokenData, err := EzspGetToken(0)
	if err != nil {
		return 0, fmt.Errorf("EzspGetToken(0) failed: %v", err)
	}

	rebootCnt = binary.LittleEndian.Uint16(tokenData)

	//rebootCnt递增并存储
	rebootCnt++
	tokenData[0] = byte(rebootCnt)
	tokenData[1] = byte(rebootCnt >> 8)
	err = EzspSetToken(0, tokenData)
	if err != nil {
		return rebootCnt, fmt.Errorf("EzspSetToken(0) failed: %v", err)
	}
	return
}

func NcpPrintAddressTable() {
	var err error
	var active bool
	var nodeID uint16
	var eui64 uint64
	for i := byte(0); i < 64; i++ {
		active, err = EzspAddressTableEntryIsActive(i)
		if err != nil {
			common.Log.Errorf("EzspAddressTableEntryIsActive(%d) failed: %v", i, err)
		}
		nodeID, err = EzspGetAddressTableRemoteNodeId(i)
		if err != nil {
			common.Log.Errorf("EzspGetAddressTableRemoteNodeId(%d) failed: %v", i, err)
		}
		eui64, err = EzspGetAddressTableRemoteEui64(i)
		if err != nil {
			common.Log.Errorf("EzspGetAddressTableRemoteEui64(%d) failed: %v", i, err)
		}
		common.Log.Infof("address table %d active=%v nodeID=0x%04x EUI64=%016x", i, active, nodeID, eui64)
	}
}

// Called when the stack status changes, usually as a result of an
// attempt to form, join, or leave a network.
func EzspStackStatusHandler(emberStatus byte) {
	switch emberStatus {
	case EMBER_NETWORK_UP, EMBER_TRUST_CENTER_EUI_HAS_CHANGED, EMBER_CHANNEL_CHANGED: // also means NETWORK_UP
		MeshStatusUp = true

		nodeType, parameters, err := EzspGetNetworkParameters()
		if err != nil {
			common.Log.Errorf("EzspGetNetworkParameters failed: %v", err)
		} else {
			MeshInfo.PANID = parameters.PanId
			MeshInfo.Channel = parameters.RadioChannel
			MeshInfo.ExPANID = parameters.ExtendedPanId

			ncpTrace("EMBER_NETWORK_UP NodeType = %d channels = %d panId = 0x%04x expanid = %016x",
				nodeType,
				parameters.RadioChannel,
				parameters.PanId,
				parameters.ExtendedPanId)
		}

	case EMBER_NETWORK_DOWN, EMBER_RECEIVED_KEY_IN_THE_CLEAR, EMBER_NO_NETWORK_KEY_RECEIVED, EMBER_NO_LINK_KEY_RECEIVED, EMBER_PRECONFIGURED_KEY_REQUIRED, EMBER_MOVE_FAILED, EMBER_JOIN_FAILED, EMBER_NO_BEACONS, EMBER_CANNOT_JOIN_AS_ROUTER:
		MeshStatusUp = false
		ncpTrace("EMBER_NETWORK_DOWN")

	default:
		common.Log.Errorf("unknown status = 0x%02x", emberStatus)
	}
}

func EzspMessageSentHandler(outgoingMessageType byte,
	indexOrDestination uint16,
	apsFrame *EmberApsFrame,
	messageTag byte,
	emberStatus byte,
	message []byte) {
	ncpTrace("%s message sent(%s) to 0x%04x, Profile 0x%04x, Cluster 0x%04x: 0x%x",
		outgoingMessageTypeToString(outgoingMessageType), emberStatusToString(emberStatus), indexOrDestination, apsFrame.ProfileId, apsFrame.ClusterId, message)
	if NcpCallbacks.NcpMessageSentHandler != nil {
		NcpCallbacks.NcpMessageSentHandler(outgoingMessageType,
			indexOrDestination,
			apsFrame,
			messageTag,
			emberStatus,
			message)
	}
}

func EzspIncomingSenderEui64Handler(senderEui64 uint64) {
	ncpTrace("Incoming sender EUI64 %016x", senderEui64)
	if NcpCallbacks.NcpIncomingSenderEui64Handler != nil {
		NcpCallbacks.NcpIncomingSenderEui64Handler(senderEui64)
	}
}

func EzspIncomingMessageHandler(incomingMessageType byte,
	apsFrame *EmberApsFrame,
	lastHopLqi byte,
	lastHopRssi int8,
	sender uint16,
	bindingIndex byte,
	addressIndex byte,
	message []byte) {

	ncpTrace("Incoming %s message from 0x%04x, Profile 0x%04x, Cluster 0x%04x: 0x%x", incomingMessageTypeToString(incomingMessageType), sender, apsFrame.ProfileId, apsFrame.ClusterId, message)
	if NcpCallbacks.NcpIncomingMessageHandler != nil {
		NcpCallbacks.NcpIncomingMessageHandler(incomingMessageType,
			apsFrame,
			lastHopLqi,
			lastHopRssi,
			sender,
			bindingIndex,
			addressIndex,
			message)
	}
}

func EzspIncomingRouteErrorHandler(emberStatus byte, target uint16) {
	ncpTrace("Incoming route error %s for 0x%04x", emberStatusToString(emberStatus), target)
	NcpSendMTORR()
}

func EzspTrustCenterJoinHandler(newNodeId uint16,
	newNodeEui64 uint64,
	deviceUpdateStatus byte,
	joinDecision byte,
	parentOfNewNode uint16) {
	ncpTrace("Trust center has 0x%04x(%016x) joined(%d)", newNodeId, newNodeEui64, deviceUpdateStatus)
	if NcpCallbacks.NcpTrustCenterJoinHandler != nil {
		NcpCallbacks.NcpTrustCenterJoinHandler(newNodeId, newNodeEui64, deviceUpdateStatus, joinDecision, parentOfNewNode)
	}
}

func NcpSendMTORR() {
	if MeshStatusUp {
		err := EzspSendManyToOneRouteRequest(EMBER_HIGH_RAM_CONCENTRATOR, 0)
		if err != nil {
			common.Log.Errorf("EzspSendManyToOneRouteRequest failed: %v", err)
		}
	}
}
