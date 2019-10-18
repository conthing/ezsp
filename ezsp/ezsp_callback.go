package ezsp

import (
	"encoding/binary"

	"github.com/conthing/utils/common"
)

func EzspTick() {
	select {
	case cb := <-callbackCh:
		ezspCallbackDispatch(cb)
	}
}

func ezspCallbackDispatch(cb *EzspFrame) {

	if cb == nil {
		common.Log.Errorf("ezspCallbackDispatch with nil frame")
		return
	}
	if cb.Data == nil {
		common.Log.Errorf("ezspCallbackDispatch with nil Data")
		return
	}

	common.Log.Debugf("callback:%s", frameIDToName(cb.FrameID))

	switch cb.FrameID {
	case EZSP_INCOMING_MESSAGE_HANDLER:
		if len(cb.Data) < 17 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}

		incomingMessageType := cb.Data[0]
		apsFrame := EmberApsFrame{}
		apsFrame.ProfileId = binary.LittleEndian.Uint16(cb.Data[1:])
		apsFrame.ClusterId = binary.LittleEndian.Uint16(cb.Data[3:])
		apsFrame.SourceEndpoint = cb.Data[5]
		apsFrame.DestinationEndpoint = cb.Data[6]
		apsFrame.Options = binary.LittleEndian.Uint16(cb.Data[7:])
		apsFrame.GroupId = binary.LittleEndian.Uint16(cb.Data[9:])
		apsFrame.Sequence = cb.Data[11]
		lastHopLqi := cb.Data[12]
		lastHopRssi := int8(cb.Data[13])
		sender := binary.LittleEndian.Uint16(cb.Data[14:])
		bindingIndex := cb.Data[16]
		addressIndex := cb.Data[17]
		messageLength := cb.Data[18]
		if len(cb.Data) != int(messageLength)+19 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}

		EzspIncomingMessageHandler(incomingMessageType,
			&apsFrame,
			lastHopLqi,
			lastHopRssi,
			sender,
			bindingIndex,
			addressIndex,
			cb.Data[19:])

	case EZSP_STACK_STATUS_HANDLER:
		if len(cb.Data) != 1 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		emberStatus := cb.Data[0]
		EzspStackStatusHandler(emberStatus)

	case EZSP_INCOMING_SENDER_EUI64_HANDLER:
		if len(cb.Data) != 8 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		senderEui64 := binary.LittleEndian.Uint64(cb.Data)
		EzspIncomingSenderEui64Handler(senderEui64)

	case EZSP_MESSAGE_SENT_HANDLER:
		if len(cb.Data) < 17 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		outgoingMessageType := cb.Data[0]
		indexOrDestination := binary.LittleEndian.Uint16(cb.Data[1:])
		apsFrame := EmberApsFrame{}
		apsFrame.ProfileId = binary.LittleEndian.Uint16(cb.Data[3:])
		apsFrame.ClusterId = binary.LittleEndian.Uint16(cb.Data[5:])
		apsFrame.SourceEndpoint = cb.Data[7]
		apsFrame.DestinationEndpoint = cb.Data[8]
		apsFrame.Options = binary.LittleEndian.Uint16(cb.Data[9:])
		apsFrame.GroupId = binary.LittleEndian.Uint16(cb.Data[11:])
		apsFrame.Sequence = cb.Data[13]
		messageTag := cb.Data[14]
		emberStatus := cb.Data[15]
		messageLength := cb.Data[16]
		if len(cb.Data) != int(messageLength)+17 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}

		EzspMessageSentHandler(outgoingMessageType,
			indexOrDestination,
			&apsFrame,
			messageTag,
			emberStatus,
			cb.Data[17:])

	case EZSP_INCOMING_ROUTE_ERROR_HANDLER:
		if len(cb.Data) != 3 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		emberStatus := cb.Data[0]
		target := binary.LittleEndian.Uint16(cb.Data[1:])
		EzspIncomingRouteErrorHandler(emberStatus, target)

	case EZSP_TRUST_CENTER_JOIN_HANDLER:
		if len(cb.Data) != 14 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		newNodeId := binary.LittleEndian.Uint16(cb.Data)
		newNodeEui64 := binary.LittleEndian.Uint64(cb.Data[2:])
		deviceUpdateStatus := cb.Data[10]
		joinDecision := cb.Data[11]
		parentOfNewNode := binary.LittleEndian.Uint16(cb.Data[12:])
		EzspTrustCenterJoinHandler(newNodeId,
			newNodeEui64,
			deviceUpdateStatus,
			joinDecision,
			parentOfNewNode)

	case EZSP_NO_CALLBACKS:
	case EZSP_STACK_TOKEN_CHANGED_HANDLER:
	case EZSP_TIMER_HANDLER:
	case EZSP_COUNTER_ROLLOVER_HANDLER:
	case EZSP_CUSTOM_FRAME_HANDLER:
	case EZSP_ENERGY_SCAN_RESULT_HANDLER:
		if len(cb.Data) != 2 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		channel := cb.Data[0]
		maxRssiValue := int8(cb.Data[1])
		EzspEnergyScanResultHandler(channel, maxRssiValue)
	case EZSP_NETWORK_FOUND_HANDLER:
		if len(cb.Data) != 16 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		networkFound := EmberZigbeeNetwork{}
		networkFound.Channel = cb.Data[0]
		networkFound.PanId = binary.LittleEndian.Uint16(cb.Data[1:])
		networkFound.ExtendedPanId = binary.LittleEndian.Uint64(cb.Data[3:])
		networkFound.AllowingJoin = (cb.Data[11] != 0)
		networkFound.StackProfile = cb.Data[12]
		networkFound.NwkUpdateId = cb.Data[13]
		lqi := cb.Data[14]
		rssi := int8(cb.Data[15])
		EzspNetworkFoundHandler(&networkFound, lqi, rssi)
	case EZSP_SCAN_COMPLETE_HANDLER:
		if len(cb.Data) != 2 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		channel := cb.Data[0]
		emberStatus := cb.Data[1]
		EzspScanCompleteHandler(channel, emberStatus)
	case EZSP_CHILD_JOIN_HANDLER:
	case EZSP_REMOTE_SET_BINDING_HANDLER:
	case EZSP_REMOTE_DELETE_BINDING_HANDLER:
	case EZSP_POLL_COMPLETE_HANDLER:
	case EZSP_POLL_HANDLER:
	case EZSP_INCOMING_ROUTE_RECORD_HANDLER:
		if len(cb.Data) < 13 {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		source := binary.LittleEndian.Uint16(cb.Data)
		sourceEui := binary.LittleEndian.Uint64(cb.Data[2:])
		lastHopLqi := cb.Data[10]
		lastHopRssi := int8(cb.Data[11])
		relayCount := cb.Data[12]
		if len(cb.Data) != 13+2*int(relayCount) {
			common.Log.Errorf("ezspCallbackDispatch %s with invalid Data length %d", frameIDToName(cb.FrameID), len(cb.Data))
			return
		}
		relay := make([]uint16, relayCount)
		for i := range relay {
			relay[i] = binary.LittleEndian.Uint16(cb.Data[13+i*2:])
		}

		EzspIncomingRouteRecordHandler(source, sourceEui, lastHopLqi, lastHopRssi, relay)
	case EZSP_INCOMING_MANY_TO_ONE_ROUTE_REQUEST_HANDLER:
	case EZSP_ID_CONFLICT_HANDLER:
	case EZSP_MAC_PASSTHROUGH_MESSAGE_HANDLER:
	case EZSP_MAC_FILTER_MATCH_MESSAGE_HANDLER:
	case EZSP_RAW_TRANSMIT_COMPLETE_HANDLER:
	case EZSP_SWITCH_NETWORK_KEY_HANDLER:
	case EZSP_ZIGBEE_KEY_ESTABLISHMENT_HANDLER:
	case EZSP_GENERATE_CBKE_KEYS_HANDLER:
	case EZSP_CALCULATE_SMACS_HANDLER:
	case EZSP_GENERATE_CBKE_KEYS_HANDLER283K1:
	case EZSP_CALCULATE_SMACS_HANDLER283K1:
	case EZSP_DSA_SIGN_HANDLER:
	case EZSP_DSA_VERIFY_HANDLER:
	case EZSP_MFGLIB_RX_HANDLER:
	case EZSP_INCOMING_BOOTLOAD_MESSAGE_HANDLER:
	case EZSP_BOOTLOAD_TRANSMIT_COMPLETE_HANDLER:
	case EZSP_ZLL_NETWORK_FOUND_HANDLER:
	case EZSP_ZLL_SCAN_COMPLETE_HANDLER:
	case EZSP_ZLL_ADDRESS_ASSIGNMENT_HANDLER:
	case EZSP_ZLL_TOUCH_LINK_TARGET_HANDLER:
	case EZSP_RF4CE_INCOMING_MESSAGE_HANDLER:
	case EZSP_RF4CE_MESSAGE_SENT_HANDLER:
	case EZSP_RF4CE_DISCOVERY_COMPLETE_HANDLER:
	case EZSP_RF4CE_DISCOVERY_REQUEST_HANDLER:
	case EZSP_RF4CE_DISCOVERY_RESPONSE_HANDLER:
	case EZSP_RF4CE_AUTO_DISCOVERY_RESPONSE_COMPLETE_HANDLER:
	case EZSP_RF4CE_PAIR_COMPLETE_HANDLER:
	case EZSP_RF4CE_PAIR_REQUEST_HANDLER:
	case EZSP_RF4CE_UNPAIR_HANDLER:
	case EZSP_RF4CE_UNPAIR_COMPLETE_HANDLER:
	default:
		common.Log.Errorf("ezspCallbackDispatch unknown callback id 0x%x", cb.FrameID)
	}
}
