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

var isCallbackIDMap [256]bool

var sequence byte
var seqMutex sync.Mutex

// callback 发送到这个ch
var callbackCh = make(chan *EzspFrame, 1)

// 用sequence做key的数组，存放收到的response时发往的ch
var responseChMap [256]chan *EzspFrame

func ezspFrameTrace(format string, v ...interface{}) {
	common.Log.Debugf(format, v...)
}

func (ezspFrame EzspFrame) String() (s string) {
	s = frameIDToName(ezspFrame.FrameID)
	if ezspFrame.Callback == 2 {
		s += "(async)"
	} else if ezspFrame.Callback == 1 {
		s += "(sync)"
	}
	s += fmt.Sprintf(" seq=0x%x %x", ezspFrame.Sequence, ezspFrame.Data)
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
// todo 应该在 AshReset 成功后再次被调用
func EzspFrameInitVariables() {
	sequence = 0

	// 清空 callbackCh
	select {
	case <-callbackCh:
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

func frameIDToName(id byte) string {
	name, ok := frameIDNameMap[id]
	if !ok {
		name = fmt.Sprintf("UNKNOWN_ID_%02X", id)
	}
	return name
}

func isValidCallbackID(callbackID byte) bool {
	if isCallbackIDMap[allCallbackIDs[0]] == false {
		for _, id := range allCallbackIDs {
			isCallbackIDMap[id] = true
		}
	}
	return isCallbackIDMap[callbackID]
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
	ezspFrame, err := ezspFrameParse(data)
	if err != nil {
		return fmt.Errorf("EZSP frame parse error: %v", err)
	}
	ezspFrameTrace("EZSP recv < %s", ezspFrame)
	if ezspFrame.Callback == 2 { // async callback 给 callbackCh
		callbackCh <- ezspFrame
		return nil
	}
	if ezspFrame.Callback == 1 { // sync callback 也给 callbackCh，另外发个nil给堵塞的发送函数
		callbackCh <- ezspFrame
	}
	ch := responseChMap[ezspFrame.Sequence]
	if ch != nil {
		if ezspFrame.Callback == 1 {
			ch <- nil
			return nil
		}
		ch <- ezspFrame
	}
	return nil
}

func EzspFrameSend(frmID byte, data []byte) (*EzspFrame, error) {
	seq := getSequence()
	ashFrm := []byte{seq, 0, frmID}
	if data != nil {
		ashFrm = append(ashFrm, data...)
	}

	// 创建接收回复的ch
	responseChMapClear(seq) //如果上一轮sequence发送时超时，有可能没有close
	responseChMap[seq] = make(chan *EzspFrame, 1)
	if responseChMap[seq] == nil {
		return nil, fmt.Errorf("EZSP send %s(seq=%d) failed: make chan failed", frameIDToName(frmID), seq)
	}

	err := ash.AshSend(ashFrm)
	if err != nil {
		responseChMapClear(seq)
		return nil, fmt.Errorf("EZSP send %s(seq=%d) failed: ash send failed: %v", frameIDToName(frmID), seq, err)
	}
	ezspFrameTrace("EZSP send > %s seq=0x%x %x", frameIDToName(frmID), seq, data)

	select {
	case response := <-responseChMap[seq]:
		close(responseChMap[seq])
		responseChMap[seq] = nil
		return response, nil
	case <-time.After(time.Millisecond * 3000):
		responseChMapClear(seq)
		return nil, fmt.Errorf("EZSP send %s timeout", frameIDToName(frmID))
	}
}

var frameIDNameMap = map[byte]string{
	// Frame ID
	// Configuration Frames
	EZSP_VERSION:                              "EZSP_VERSION",
	EZSP_GET_CONFIGURATION_VALUE:              "EZSP_GET_CONFIGURATION_VALUE",
	EZSP_SET_CONFIGURATION_VALUE:              "EZSP_SET_CONFIGURATION_VALUE",
	EZSP_ADD_ENDPOINT:                         "EZSP_ADD_ENDPOINT",
	EZSP_SET_POLICY:                           "EZSP_SET_POLICY",
	EZSP_GET_POLICY:                           "EZSP_GET_POLICY",
	EZSP_GET_VALUE:                            "EZSP_GET_VALUE",
	EZSP_GET_EXTENDED_VALUE:                   "EZSP_GET_EXTENDED_VALUE",
	EZSP_SET_VALUE:                            "EZSP_SET_VALUE",
	EZSP_SET_GPIO_CURRENT_CONFIGURATION:       "EZSP_SET_GPIO_CURRENT_CONFIGURATION",
	EZSP_SET_GPIO_POWER_UP_DOWN_CONFIGURATION: "EZSP_SET_GPIO_POWER_UP_DOWN_CONFIGURATION",
	EZSP_SET_GPIO_RADIO_POWER_MASK:            "EZSP_SET_GPIO_RADIO_POWER_MASK",

	// Utilities Frames
	EZSP_NOP:                         "EZSP_NOP",
	EZSP_ECHO:                        "EZSP_ECHO",
	EZSP_INVALID_COMMAND:             "EZSP_INVALID_COMMAND",
	EZSP_CALLBACK:                    "EZSP_CALLBACK",
	EZSP_NO_CALLBACKS:                "EZSP_NO_CALLBACKS",
	EZSP_SET_TOKEN:                   "EZSP_SET_TOKEN",
	EZSP_GET_TOKEN:                   "EZSP_GET_TOKEN",
	EZSP_GET_MFG_TOKEN:               "EZSP_GET_MFG_TOKEN",
	EZSP_SET_MFG_TOKEN:               "EZSP_SET_MFG_TOKEN",
	EZSP_STACK_TOKEN_CHANGED_HANDLER: "EZSP_STACK_TOKEN_CHANGED_HANDLER",
	EZSP_GET_RANDOM_NUMBER:           "EZSP_GET_RANDOM_NUMBER",
	EZSP_SET_TIMER:                   "EZSP_SET_TIMER",
	EZSP_GET_TIMER:                   "EZSP_GET_TIMER",
	EZSP_TIMER_HANDLER:               "EZSP_TIMER_HANDLER",
	EZSP_DEBUG_WRITE:                 "EZSP_DEBUG_WRITE",
	EZSP_READ_AND_CLEAR_COUNTERS:     "EZSP_READ_AND_CLEAR_COUNTERS",
	EZSP_READ_COUNTERS:               "EZSP_READ_COUNTERS",
	EZSP_COUNTER_ROLLOVER_HANDLER:    "EZSP_COUNTER_ROLLOVER_HANDLER",
	EZSP_DELAY_TEST:                  "EZSP_DELAY_TEST",
	EZSP_GET_LIBRARY_STATUS:          "EZSP_GET_LIBRARY_STATUS",
	EZSP_GET_XNCP_INFO:               "EZSP_GET_XNCP_INFO",
	EZSP_CUSTOM_FRAME:                "EZSP_CUSTOM_FRAME",
	EZSP_CUSTOM_FRAME_HANDLER:        "EZSP_CUSTOM_FRAME_HANDLER",

	// Networking Frames
	EZSP_SET_MANUFACTURER_CODE:       "EZSP_SET_MANUFACTURER_CODE",
	EZSP_SET_POWER_DESCRIPTOR:        "EZSP_SET_POWER_DESCRIPTOR",
	EZSP_NETWORK_INIT:                "EZSP_NETWORK_INIT",
	EZSP_NETWORK_INIT_EXTENDED:       "EZSP_NETWORK_INIT_EXTENDED",
	EZSP_NETWORK_STATE:               "EZSP_NETWORK_STATE",
	EZSP_STACK_STATUS_HANDLER:        "EZSP_STACK_STATUS_HANDLER",
	EZSP_START_SCAN:                  "EZSP_START_SCAN",
	EZSP_ENERGY_SCAN_RESULT_HANDLER:  "EZSP_ENERGY_SCAN_RESULT_HANDLER",
	EZSP_NETWORK_FOUND_HANDLER:       "EZSP_NETWORK_FOUND_HANDLER",
	EZSP_SCAN_COMPLETE_HANDLER:       "EZSP_SCAN_COMPLETE_HANDLER",
	EZSP_STOP_SCAN:                   "EZSP_STOP_SCAN",
	EZSP_FORM_NETWORK:                "EZSP_FORM_NETWORK",
	EZSP_JOIN_NETWORK:                "EZSP_JOIN_NETWORK",
	EZSP_LEAVE_NETWORK:               "EZSP_LEAVE_NETWORK",
	EZSP_FIND_AND_REJOIN_NETWORK:     "EZSP_FIND_AND_REJOIN_NETWORK",
	EZSP_PERMIT_JOINING:              "EZSP_PERMIT_JOINING",
	EZSP_CHILD_JOIN_HANDLER:          "EZSP_CHILD_JOIN_HANDLER",
	EZSP_ENERGY_SCAN_REQUEST:         "EZSP_ENERGY_SCAN_REQUEST",
	EZSP_GET_EUI64:                   "EZSP_GET_EUI64",
	EZSP_GET_NODE_ID:                 "EZSP_GET_NODE_ID",
	EZSP_GET_NETWORK_PARAMETERS:      "EZSP_GET_NETWORK_PARAMETERS",
	EZSP_GET_PARENT_CHILD_PARAMETERS: "EZSP_GET_PARENT_CHILD_PARAMETERS",
	EZSP_GET_CHILD_DATA:              "EZSP_GET_CHILD_DATA",
	EZSP_GET_NEIGHBOR:                "EZSP_GET_NEIGHBOR",
	EZSP_NEIGHBOR_COUNT:              "EZSP_NEIGHBOR_COUNT",
	EZSP_GET_ROUTE_TABLE_ENTRY:       "EZSP_GET_ROUTE_TABLE_ENTRY",
	EZSP_SET_RADIO_POWER:             "EZSP_SET_RADIO_POWER",
	EZSP_SET_RADIO_CHANNEL:           "EZSP_SET_RADIO_CHANNEL",
	EZSP_SET_CONCENTRATOR:            "EZSP_SET_CONCENTRATOR",

	// Binding Frames
	EZSP_CLEAR_BINDING_TABLE:           "EZSP_CLEAR_BINDING_TABLE",
	EZSP_SET_BINDING:                   "EZSP_SET_BINDING",
	EZSP_GET_BINDING:                   "EZSP_GET_BINDING",
	EZSP_DELETE_BINDING:                "EZSP_DELETE_BINDING",
	EZSP_BINDING_IS_ACTIVE:             "EZSP_BINDING_IS_ACTIVE",
	EZSP_GET_BINDING_REMOTE_NODE_ID:    "EZSP_GET_BINDING_REMOTE_NODE_ID",
	EZSP_SET_BINDING_REMOTE_NODE_ID:    "EZSP_SET_BINDING_REMOTE_NODE_ID",
	EZSP_REMOTE_SET_BINDING_HANDLER:    "EZSP_REMOTE_SET_BINDING_HANDLER",
	EZSP_REMOTE_DELETE_BINDING_HANDLER: "EZSP_REMOTE_DELETE_BINDING_HANDLER",

	// Messaging Frames
	EZSP_MAXIMUM_PAYLOAD_LENGTH:                     "EZSP_MAXIMUM_PAYLOAD_LENGTH",
	EZSP_SEND_UNICAST:                               "EZSP_SEND_UNICAST",
	EZSP_SEND_BROADCAST:                             "EZSP_SEND_BROADCAST",
	EZSP_PROXY_BROADCAST:                            "EZSP_PROXY_BROADCAST",
	EZSP_SEND_MULTICAST:                             "EZSP_SEND_MULTICAST",
	EZSP_SEND_REPLY:                                 "EZSP_SEND_REPLY",
	EZSP_MESSAGE_SENT_HANDLER:                       "EZSP_MESSAGE_SENT_HANDLER",
	EZSP_SEND_MANY_TO_ONE_ROUTE_REQUEST:             "EZSP_SEND_MANY_TO_ONE_ROUTE_REQUEST",
	EZSP_POLL_FOR_DATA:                              "EZSP_POLL_FOR_DATA",
	EZSP_POLL_COMPLETE_HANDLER:                      "EZSP_POLL_COMPLETE_HANDLER",
	EZSP_POLL_HANDLER:                               "EZSP_POLL_HANDLER",
	EZSP_INCOMING_SENDER_EUI64_HANDLER:              "EZSP_INCOMING_SENDER_EUI64_HANDLER",
	EZSP_INCOMING_MESSAGE_HANDLER:                   "EZSP_INCOMING_MESSAGE_HANDLER",
	EZSP_INCOMING_ROUTE_RECORD_HANDLER:              "EZSP_INCOMING_ROUTE_RECORD_HANDLER",
	EZSP_SET_SOURCE_ROUTE:                           "EZSP_SET_SOURCE_ROUTE",
	EZSP_INCOMING_MANY_TO_ONE_ROUTE_REQUEST_HANDLER: "EZSP_INCOMING_MANY_TO_ONE_ROUTE_REQUEST_HANDLER",
	EZSP_INCOMING_ROUTE_ERROR_HANDLER:               "EZSP_INCOMING_ROUTE_ERROR_HANDLER",
	EZSP_ADDRESS_TABLE_ENTRY_IS_ACTIVE:              "EZSP_ADDRESS_TABLE_ENTRY_IS_ACTIVE",
	EZSP_SET_ADDRESS_TABLE_REMOTE_EUI64:             "EZSP_SET_ADDRESS_TABLE_REMOTE_EUI64",
	EZSP_SET_ADDRESS_TABLE_REMOTE_NODE_ID:           "EZSP_SET_ADDRESS_TABLE_REMOTE_NODE_ID",
	EZSP_GET_ADDRESS_TABLE_REMOTE_EUI64:             "EZSP_GET_ADDRESS_TABLE_REMOTE_EUI64",
	EZSP_GET_ADDRESS_TABLE_REMOTE_NODE_ID:           "EZSP_GET_ADDRESS_TABLE_REMOTE_NODE_ID",
	EZSP_SET_EXTENDED_TIMEOUT:                       "EZSP_SET_EXTENDED_TIMEOUT",
	EZSP_GET_EXTENDED_TIMEOUT:                       "EZSP_GET_EXTENDED_TIMEOUT",
	EZSP_REPLACE_ADDRESS_TABLE_ENTRY:                "EZSP_REPLACE_ADDRESS_TABLE_ENTRY",
	EZSP_LOOKUP_NODE_ID_BY_EUI64:                    "EZSP_LOOKUP_NODE_ID_BY_EUI64",
	EZSP_LOOKUP_EUI64_BY_NODE_ID:                    "EZSP_LOOKUP_EUI64_BY_NODE_ID",
	EZSP_GET_MULTICAST_TABLE_ENTRY:                  "EZSP_GET_MULTICAST_TABLE_ENTRY",
	EZSP_SET_MULTICAST_TABLE_ENTRY:                  "EZSP_SET_MULTICAST_TABLE_ENTRY",
	EZSP_ID_CONFLICT_HANDLER:                        "EZSP_ID_CONFLICT_HANDLER",
	EZSP_SEND_RAW_MESSAGE:                           "EZSP_SEND_RAW_MESSAGE",
	EZSP_MAC_PASSTHROUGH_MESSAGE_HANDLER:            "EZSP_MAC_PASSTHROUGH_MESSAGE_HANDLER",
	EZSP_MAC_FILTER_MATCH_MESSAGE_HANDLER:           "EZSP_MAC_FILTER_MATCH_MESSAGE_HANDLER",
	EZSP_RAW_TRANSMIT_COMPLETE_HANDLER:              "EZSP_RAW_TRANSMIT_COMPLETE_HANDLER",

	// Security Frames
	EZSP_SET_INITIAL_SECURITY_STATE:       "EZSP_SET_INITIAL_SECURITY_STATE",
	EZSP_GET_CURRENT_SECURITY_STATE:       "EZSP_GET_CURRENT_SECURITY_STATE",
	EZSP_GET_KEY:                          "EZSP_GET_KEY",
	EZSP_SWITCH_NETWORK_KEY_HANDLER:       "EZSP_SWITCH_NETWORK_KEY_HANDLER",
	EZSP_GET_KEY_TABLE_ENTRY:              "EZSP_GET_KEY_TABLE_ENTRY",
	EZSP_SET_KEY_TABLE_ENTRY:              "EZSP_SET_KEY_TABLE_ENTRY",
	EZSP_FIND_KEY_TABLE_ENTRY:             "EZSP_FIND_KEY_TABLE_ENTRY",
	EZSP_ADD_OR_UPDATE_KEY_TABLE_ENTRY:    "EZSP_ADD_OR_UPDATE_KEY_TABLE_ENTRY",
	EZSP_ERASE_KEY_TABLE_ENTRY:            "EZSP_ERASE_KEY_TABLE_ENTRY",
	EZSP_CLEAR_KEY_TABLE:                  "EZSP_CLEAR_KEY_TABLE",
	EZSP_REQUEST_LINK_KEY:                 "EZSP_REQUEST_LINK_KEY",
	EZSP_ZIGBEE_KEY_ESTABLISHMENT_HANDLER: "EZSP_ZIGBEE_KEY_ESTABLISHMENT_HANDLER",

	// Trust Center Frames
	EZSP_TRUST_CENTER_JOIN_HANDLER:    "EZSP_TRUST_CENTER_JOIN_HANDLER",
	EZSP_BROADCAST_NEXT_NETWORK_KEY:   "EZSP_BROADCAST_NEXT_NETWORK_KEY",
	EZSP_BROADCAST_NETWORK_KEY_SWITCH: "EZSP_BROADCAST_NETWORK_KEY_SWITCH",
	EZSP_BECOME_TRUST_CENTER:          "EZSP_BECOME_TRUST_CENTER",
	EZSP_AES_MMO_HASH:                 "EZSP_AES_MMO_HASH",
	EZSP_REMOVE_DEVICE:                "EZSP_REMOVE_DEVICE",
	EZSP_UNICAST_NWK_KEY_UPDATE:       "EZSP_UNICAST_NWK_KEY_UPDATE",

	// Certificate Based Key Exchange (CBKE(
	EZSP_GENERATE_CBKE_KEYS:                             "EZSP_GENERATE_CBKE_KEYS",
	EZSP_GENERATE_CBKE_KEYS_HANDLER:                     "EZSP_GENERATE_CBKE_KEYS_HANDLER",
	EZSP_CALCULATE_SMACS:                                "EZSP_CALCULATE_SMACS",
	EZSP_CALCULATE_SMACS_HANDLER:                        "EZSP_CALCULATE_SMACS_HANDLER",
	EZSP_GENERATE_CBKE_KEYS283K1:                        "EZSP_GENERATE_CBKE_KEYS283K1",
	EZSP_GENERATE_CBKE_KEYS_HANDLER283K1:                "EZSP_GENERATE_CBKE_KEYS_HANDLER283K1",
	EZSP_CALCULATE_SMACS283K1:                           "EZSP_CALCULATE_SMACS283K1",
	EZSP_CALCULATE_SMACS_HANDLER283K1:                   "EZSP_CALCULATE_SMACS_HANDLER283K1",
	EZSP_CLEAR_TEMPORARY_DATA_MAYBE_STORE_LINK_KEY:      "EZSP_CLEAR_TEMPORARY_DATA_MAYBE_STORE_LINK_KEY",
	EZSP_CLEAR_TEMPORARY_DATA_MAYBE_STORE_LINK_KEY283K1: "EZSP_CLEAR_TEMPORARY_DATA_MAYBE_STORE_LINK_KEY283K1",
	EZSP_GET_CERTIFICATE:                                "EZSP_GET_CERTIFICATE",
	EZSP_GET_CERTIFICATE283K1:                           "EZSP_GET_CERTIFICATE283K1",
	EZSP_DSA_SIGN:                                       "EZSP_DSA_SIGN",
	EZSP_DSA_SIGN_HANDLER:                               "EZSP_DSA_SIGN_HANDLER",
	EZSP_DSA_VERIFY:                                     "EZSP_DSA_VERIFY",
	EZSP_DSA_VERIFY_HANDLER:                             "EZSP_DSA_VERIFY_HANDLER",
	EZSP_SET_PREINSTALLED_CBKE_DATA:                     "EZSP_SET_PREINSTALLED_CBKE_DATA",
	EZSP_SAVE_PREINSTALLED_CBKE_DATA283K1:               "EZSP_SAVE_PREINSTALLED_CBKE_DATA283K1",

	// Mfglib
	EZSP_MFGLIB_START:        "EZSP_MFGLIB_START",
	EZSP_MFGLIB_END:          "EZSP_MFGLIB_END",
	EZSP_MFGLIB_START_TONE:   "EZSP_MFGLIB_START_TONE",
	EZSP_MFGLIB_STOP_TONE:    "EZSP_MFGLIB_STOP_TONE",
	EZSP_MFGLIB_START_STREAM: "EZSP_MFGLIB_START_STREAM",
	EZSP_MFGLIB_STOP_STREAM:  "EZSP_MFGLIB_STOP_STREAM",
	EZSP_MFGLIB_SEND_PACKET:  "EZSP_MFGLIB_SEND_PACKET",
	EZSP_MFGLIB_SET_CHANNEL:  "EZSP_MFGLIB_SET_CHANNEL",
	EZSP_MFGLIB_GET_CHANNEL:  "EZSP_MFGLIB_GET_CHANNEL",
	EZSP_MFGLIB_SET_POWER:    "EZSP_MFGLIB_SET_POWER",
	EZSP_MFGLIB_GET_POWER:    "EZSP_MFGLIB_GET_POWER",
	EZSP_MFGLIB_RX_HANDLER:   "EZSP_MFGLIB_RX_HANDLER",

	// Bootloader
	EZSP_LAUNCH_STANDALONE_BOOTLOADER:                     "EZSP_LAUNCH_STANDALONE_BOOTLOADER",
	EZSP_SEND_BOOTLOAD_MESSAGE:                            "EZSP_SEND_BOOTLOAD_MESSAGE",
	EZSP_GET_STANDALONE_BOOTLOADER_VERSION_PLAT_MICRO_PHY: "EZSP_GET_STANDALONE_BOOTLOADER_VERSION_PLAT_MICRO_PHY",
	EZSP_INCOMING_BOOTLOAD_MESSAGE_HANDLER:                "EZSP_INCOMING_BOOTLOAD_MESSAGE_HANDLER",
	EZSP_BOOTLOAD_TRANSMIT_COMPLETE_HANDLER:               "EZSP_BOOTLOAD_TRANSMIT_COMPLETE_HANDLER",
	EZSP_AES_ENCRYPT:                                      "EZSP_AES_ENCRYPT",
	EZSP_OVERRIDE_CURRENT_CHANNEL:                         "EZSP_OVERRIDE_CURRENT_CHANNEL",

	// ZLL
	EZSP_ZLL_NETWORK_OPS:                "EZSP_ZLL_NETWORK_OPS",
	EZSP_ZLL_SET_INITIAL_SECURITY_STATE: "EZSP_ZLL_SET_INITIAL_SECURITY_STATE",
	EZSP_ZLL_START_SCAN:                 "EZSP_ZLL_START_SCAN",
	EZSP_ZLL_SET_RX_ON_WHEN_IDLE:        "EZSP_ZLL_SET_RX_ON_WHEN_IDLE",
	EZSP_ZLL_NETWORK_FOUND_HANDLER:      "EZSP_ZLL_NETWORK_FOUND_HANDLER",
	EZSP_ZLL_SCAN_COMPLETE_HANDLER:      "EZSP_ZLL_SCAN_COMPLETE_HANDLER",
	EZSP_ZLL_ADDRESS_ASSIGNMENT_HANDLER: "EZSP_ZLL_ADDRESS_ASSIGNMENT_HANDLER",
	EZSP_SET_LOGICAL_AND_RADIO_CHANNEL:  "EZSP_SET_LOGICAL_AND_RADIO_CHANNEL",
	EZSP_GET_LOGICAL_CHANNEL:            "EZSP_GET_LOGICAL_CHANNEL",
	EZSP_ZLL_TOUCH_LINK_TARGET_HANDLER:  "EZSP_ZLL_TOUCH_LINK_TARGET_HANDLER",
	EZSP_ZLL_GET_TOKENS:                 "EZSP_ZLL_GET_TOKENS",
	EZSP_ZLL_SET_DATA_TOKEN:             "EZSP_ZLL_SET_DATA_TOKEN",
	EZSP_ZLL_SET_NON_ZLL_NETWORK:        "EZSP_ZLL_SET_NON_ZLL_NETWORK",
	EZSP_IS_ZLL_NETWORK:                 "EZSP_IS_ZLL_NETWORK",

	// RF4CE
	EZSP_RF4CE_SET_PAIRING_TABLE_ENTRY:                  "EZSP_RF4CE_SET_PAIRING_TABLE_ENTRY",
	EZSP_RF4CE_GET_PAIRING_TABLE_ENTRY:                  "EZSP_RF4CE_GET_PAIRING_TABLE_ENTRY",
	EZSP_RF4CE_DELETE_PAIRING_TABLE_ENTRY:               "EZSP_RF4CE_DELETE_PAIRING_TABLE_ENTRY",
	EZSP_RF4CE_KEY_UPDATE:                               "EZSP_RF4CE_KEY_UPDATE",
	EZSP_RF4CE_SEND:                                     "EZSP_RF4CE_SEND",
	EZSP_RF4CE_INCOMING_MESSAGE_HANDLER:                 "EZSP_RF4CE_INCOMING_MESSAGE_HANDLER",
	EZSP_RF4CE_MESSAGE_SENT_HANDLER:                     "EZSP_RF4CE_MESSAGE_SENT_HANDLER",
	EZSP_RF4CE_START:                                    "EZSP_RF4CE_START",
	EZSP_RF4CE_STOP:                                     "EZSP_RF4CE_STOP",
	EZSP_RF4CE_DISCOVERY:                                "EZSP_RF4CE_DISCOVERY",
	EZSP_RF4CE_DISCOVERY_COMPLETE_HANDLER:               "EZSP_RF4CE_DISCOVERY_COMPLETE_HANDLER",
	EZSP_RF4CE_DISCOVERY_REQUEST_HANDLER:                "EZSP_RF4CE_DISCOVERY_REQUEST_HANDLER",
	EZSP_RF4CE_DISCOVERY_RESPONSE_HANDLER:               "EZSP_RF4CE_DISCOVERY_RESPONSE_HANDLER",
	EZSP_RF4CE_ENABLE_AUTO_DISCOVERY_RESPONSE:           "EZSP_RF4CE_ENABLE_AUTO_DISCOVERY_RESPONSE",
	EZSP_RF4CE_AUTO_DISCOVERY_RESPONSE_COMPLETE_HANDLER: "EZSP_RF4CE_AUTO_DISCOVERY_RESPONSE_COMPLETE_HANDLER",
	EZSP_RF4CE_PAIR:                                     "EZSP_RF4CE_PAIR",
	EZSP_RF4CE_PAIR_COMPLETE_HANDLER:                    "EZSP_RF4CE_PAIR_COMPLETE_HANDLER",
	EZSP_RF4CE_PAIR_REQUEST_HANDLER:                     "EZSP_RF4CE_PAIR_REQUEST_HANDLER",
	EZSP_RF4CE_UNPAIR:                                   "EZSP_RF4CE_UNPAIR",
	EZSP_RF4CE_UNPAIR_HANDLER:                           "EZSP_RF4CE_UNPAIR_HANDLER",
	EZSP_RF4CE_UNPAIR_COMPLETE_HANDLER:                  "EZSP_RF4CE_UNPAIR_COMPLETE_HANDLER",
	EZSP_RF4CE_SET_POWER_SAVING_PARAMETERS:              "EZSP_RF4CE_SET_POWER_SAVING_PARAMETERS",
	EZSP_RF4CE_SET_FREQUENCY_AGILITY_PARAMETERS:         "EZSP_RF4CE_SET_FREQUENCY_AGILITY_PARAMETERS",
	EZSP_RF4CE_SET_APPLICATION_INFO:                     "EZSP_RF4CE_SET_APPLICATION_INFO",
	EZSP_RF4CE_GET_APPLICATION_INFO:                     "EZSP_RF4CE_GET_APPLICATION_INFO",
	EZSP_RF4CE_GET_MAX_PAYLOAD:                          "EZSP_RF4CE_GET_MAX_PAYLOAD",
}

var allCallbackIDs = [...]byte{
	EZSP_NO_CALLBACKS,
	EZSP_STACK_TOKEN_CHANGED_HANDLER,
	EZSP_TIMER_HANDLER,
	EZSP_COUNTER_ROLLOVER_HANDLER,
	EZSP_CUSTOM_FRAME_HANDLER,
	EZSP_STACK_STATUS_HANDLER,
	EZSP_ENERGY_SCAN_RESULT_HANDLER,
	EZSP_NETWORK_FOUND_HANDLER,
	EZSP_SCAN_COMPLETE_HANDLER,
	EZSP_CHILD_JOIN_HANDLER,
	EZSP_REMOTE_SET_BINDING_HANDLER,
	EZSP_REMOTE_DELETE_BINDING_HANDLER,
	EZSP_MESSAGE_SENT_HANDLER,
	EZSP_POLL_COMPLETE_HANDLER,
	EZSP_POLL_HANDLER,
	EZSP_INCOMING_SENDER_EUI64_HANDLER,
	EZSP_INCOMING_MESSAGE_HANDLER,
	EZSP_INCOMING_ROUTE_RECORD_HANDLER,
	EZSP_INCOMING_MANY_TO_ONE_ROUTE_REQUEST_HANDLER,
	EZSP_INCOMING_ROUTE_ERROR_HANDLER,
	EZSP_ID_CONFLICT_HANDLER,
	EZSP_MAC_PASSTHROUGH_MESSAGE_HANDLER,
	EZSP_MAC_FILTER_MATCH_MESSAGE_HANDLER,
	EZSP_RAW_TRANSMIT_COMPLETE_HANDLER,
	EZSP_SWITCH_NETWORK_KEY_HANDLER,
	EZSP_ZIGBEE_KEY_ESTABLISHMENT_HANDLER,
	EZSP_TRUST_CENTER_JOIN_HANDLER,
	EZSP_GENERATE_CBKE_KEYS_HANDLER,
	EZSP_CALCULATE_SMACS_HANDLER,
	EZSP_GENERATE_CBKE_KEYS_HANDLER283K1,
	EZSP_CALCULATE_SMACS_HANDLER283K1,
	EZSP_DSA_SIGN_HANDLER,
	EZSP_DSA_VERIFY_HANDLER,
	EZSP_MFGLIB_RX_HANDLER,
	EZSP_INCOMING_BOOTLOAD_MESSAGE_HANDLER,
	EZSP_BOOTLOAD_TRANSMIT_COMPLETE_HANDLER,
	EZSP_ZLL_NETWORK_FOUND_HANDLER,
	EZSP_ZLL_SCAN_COMPLETE_HANDLER,
	EZSP_ZLL_ADDRESS_ASSIGNMENT_HANDLER,
	EZSP_ZLL_TOUCH_LINK_TARGET_HANDLER,
	EZSP_RF4CE_INCOMING_MESSAGE_HANDLER,
	EZSP_RF4CE_MESSAGE_SENT_HANDLER,
	EZSP_RF4CE_DISCOVERY_COMPLETE_HANDLER,
	EZSP_RF4CE_DISCOVERY_REQUEST_HANDLER,
	EZSP_RF4CE_DISCOVERY_RESPONSE_HANDLER,
	EZSP_RF4CE_AUTO_DISCOVERY_RESPONSE_COMPLETE_HANDLER,
	EZSP_RF4CE_PAIR_COMPLETE_HANDLER,
	EZSP_RF4CE_PAIR_REQUEST_HANDLER,
	EZSP_RF4CE_UNPAIR_HANDLER,
	EZSP_RF4CE_UNPAIR_COMPLETE_HANDLER,
}

const (
	// Frame ID
	// Configuration Frames
	EZSP_VERSION                              = byte(0x00)
	EZSP_GET_CONFIGURATION_VALUE              = byte(0x52)
	EZSP_SET_CONFIGURATION_VALUE              = byte(0x53)
	EZSP_ADD_ENDPOINT                         = byte(0x02)
	EZSP_SET_POLICY                           = byte(0x55)
	EZSP_GET_POLICY                           = byte(0x56)
	EZSP_GET_VALUE                            = byte(0xAA)
	EZSP_GET_EXTENDED_VALUE                   = byte(0x03)
	EZSP_SET_VALUE                            = byte(0xAB)
	EZSP_SET_GPIO_CURRENT_CONFIGURATION       = byte(0xAC)
	EZSP_SET_GPIO_POWER_UP_DOWN_CONFIGURATION = byte(0xAD)
	EZSP_SET_GPIO_RADIO_POWER_MASK            = byte(0xAE)

	// Utilities Frames
	EZSP_NOP                         = byte(0x05)
	EZSP_ECHO                        = byte(0x81)
	EZSP_INVALID_COMMAND             = byte(0x58)
	EZSP_CALLBACK                    = byte(0x06)
	EZSP_NO_CALLBACKS                = byte(0x07)
	EZSP_SET_TOKEN                   = byte(0x09)
	EZSP_GET_TOKEN                   = byte(0x0A)
	EZSP_GET_MFG_TOKEN               = byte(0x0B)
	EZSP_SET_MFG_TOKEN               = byte(0x0C)
	EZSP_STACK_TOKEN_CHANGED_HANDLER = byte(0x0D)
	EZSP_GET_RANDOM_NUMBER           = byte(0x49)
	EZSP_SET_TIMER                   = byte(0x0E)
	EZSP_GET_TIMER                   = byte(0x4E)
	EZSP_TIMER_HANDLER               = byte(0x0F)
	EZSP_DEBUG_WRITE                 = byte(0x12)
	EZSP_READ_AND_CLEAR_COUNTERS     = byte(0x65)
	EZSP_READ_COUNTERS               = byte(0xF1)
	EZSP_COUNTER_ROLLOVER_HANDLER    = byte(0xF2)
	EZSP_DELAY_TEST                  = byte(0x9D)
	EZSP_GET_LIBRARY_STATUS          = byte(0x01)
	EZSP_GET_XNCP_INFO               = byte(0x13)
	EZSP_CUSTOM_FRAME                = byte(0x47)
	EZSP_CUSTOM_FRAME_HANDLER        = byte(0x54)

	// Networking Frames
	EZSP_SET_MANUFACTURER_CODE       = byte(0x15)
	EZSP_SET_POWER_DESCRIPTOR        = byte(0x16)
	EZSP_NETWORK_INIT                = byte(0x17)
	EZSP_NETWORK_INIT_EXTENDED       = byte(0x70)
	EZSP_NETWORK_STATE               = byte(0x18)
	EZSP_STACK_STATUS_HANDLER        = byte(0x19)
	EZSP_START_SCAN                  = byte(0x1A)
	EZSP_ENERGY_SCAN_RESULT_HANDLER  = byte(0x48)
	EZSP_NETWORK_FOUND_HANDLER       = byte(0x1B)
	EZSP_SCAN_COMPLETE_HANDLER       = byte(0x1C)
	EZSP_STOP_SCAN                   = byte(0x1D)
	EZSP_FORM_NETWORK                = byte(0x1E)
	EZSP_JOIN_NETWORK                = byte(0x1F)
	EZSP_LEAVE_NETWORK               = byte(0x20)
	EZSP_FIND_AND_REJOIN_NETWORK     = byte(0x21)
	EZSP_PERMIT_JOINING              = byte(0x22)
	EZSP_CHILD_JOIN_HANDLER          = byte(0x23)
	EZSP_ENERGY_SCAN_REQUEST         = byte(0x9C)
	EZSP_GET_EUI64                   = byte(0x26)
	EZSP_GET_NODE_ID                 = byte(0x27)
	EZSP_GET_NETWORK_PARAMETERS      = byte(0x28)
	EZSP_GET_PARENT_CHILD_PARAMETERS = byte(0x29)
	EZSP_GET_CHILD_DATA              = byte(0x4A)
	EZSP_GET_NEIGHBOR                = byte(0x79)
	EZSP_NEIGHBOR_COUNT              = byte(0x7A)
	EZSP_GET_ROUTE_TABLE_ENTRY       = byte(0x7B)
	EZSP_SET_RADIO_POWER             = byte(0x99)
	EZSP_SET_RADIO_CHANNEL           = byte(0x9A)
	EZSP_SET_CONCENTRATOR            = byte(0x10)

	// Binding Frames
	EZSP_CLEAR_BINDING_TABLE           = byte(0x2A)
	EZSP_SET_BINDING                   = byte(0x2B)
	EZSP_GET_BINDING                   = byte(0x2C)
	EZSP_DELETE_BINDING                = byte(0x2D)
	EZSP_BINDING_IS_ACTIVE             = byte(0x2E)
	EZSP_GET_BINDING_REMOTE_NODE_ID    = byte(0x2F)
	EZSP_SET_BINDING_REMOTE_NODE_ID    = byte(0x30)
	EZSP_REMOTE_SET_BINDING_HANDLER    = byte(0x31)
	EZSP_REMOTE_DELETE_BINDING_HANDLER = byte(0x32)

	// Messaging Frames
	EZSP_MAXIMUM_PAYLOAD_LENGTH                     = byte(0x33)
	EZSP_SEND_UNICAST                               = byte(0x34)
	EZSP_SEND_BROADCAST                             = byte(0x36)
	EZSP_PROXY_BROADCAST                            = byte(0x37)
	EZSP_SEND_MULTICAST                             = byte(0x38)
	EZSP_SEND_REPLY                                 = byte(0x39)
	EZSP_MESSAGE_SENT_HANDLER                       = byte(0x3F)
	EZSP_SEND_MANY_TO_ONE_ROUTE_REQUEST             = byte(0x41)
	EZSP_POLL_FOR_DATA                              = byte(0x42)
	EZSP_POLL_COMPLETE_HANDLER                      = byte(0x43)
	EZSP_POLL_HANDLER                               = byte(0x44)
	EZSP_INCOMING_SENDER_EUI64_HANDLER              = byte(0x62)
	EZSP_INCOMING_MESSAGE_HANDLER                   = byte(0x45)
	EZSP_INCOMING_ROUTE_RECORD_HANDLER              = byte(0x59)
	EZSP_SET_SOURCE_ROUTE                           = byte(0x5A)
	EZSP_INCOMING_MANY_TO_ONE_ROUTE_REQUEST_HANDLER = byte(0x7D)
	EZSP_INCOMING_ROUTE_ERROR_HANDLER               = byte(0x80)
	EZSP_ADDRESS_TABLE_ENTRY_IS_ACTIVE              = byte(0x5B)
	EZSP_SET_ADDRESS_TABLE_REMOTE_EUI64             = byte(0x5C)
	EZSP_SET_ADDRESS_TABLE_REMOTE_NODE_ID           = byte(0x5D)
	EZSP_GET_ADDRESS_TABLE_REMOTE_EUI64             = byte(0x5E)
	EZSP_GET_ADDRESS_TABLE_REMOTE_NODE_ID           = byte(0x5F)
	EZSP_SET_EXTENDED_TIMEOUT                       = byte(0x7E)
	EZSP_GET_EXTENDED_TIMEOUT                       = byte(0x7F)
	EZSP_REPLACE_ADDRESS_TABLE_ENTRY                = byte(0x82)
	EZSP_LOOKUP_NODE_ID_BY_EUI64                    = byte(0x60)
	EZSP_LOOKUP_EUI64_BY_NODE_ID                    = byte(0x61)
	EZSP_GET_MULTICAST_TABLE_ENTRY                  = byte(0x63)
	EZSP_SET_MULTICAST_TABLE_ENTRY                  = byte(0x64)
	EZSP_ID_CONFLICT_HANDLER                        = byte(0x7C)
	EZSP_SEND_RAW_MESSAGE                           = byte(0x96)
	EZSP_MAC_PASSTHROUGH_MESSAGE_HANDLER            = byte(0x97)
	EZSP_MAC_FILTER_MATCH_MESSAGE_HANDLER           = byte(0x46)
	EZSP_RAW_TRANSMIT_COMPLETE_HANDLER              = byte(0x98)

	// Security Frames
	EZSP_SET_INITIAL_SECURITY_STATE       = byte(0x68)
	EZSP_GET_CURRENT_SECURITY_STATE       = byte(0x69)
	EZSP_GET_KEY                          = byte(0x6a)
	EZSP_SWITCH_NETWORK_KEY_HANDLER       = byte(0x6e)
	EZSP_GET_KEY_TABLE_ENTRY              = byte(0x71)
	EZSP_SET_KEY_TABLE_ENTRY              = byte(0x72)
	EZSP_FIND_KEY_TABLE_ENTRY             = byte(0x75)
	EZSP_ADD_OR_UPDATE_KEY_TABLE_ENTRY    = byte(0x66)
	EZSP_ERASE_KEY_TABLE_ENTRY            = byte(0x76)
	EZSP_CLEAR_KEY_TABLE                  = byte(0xB1)
	EZSP_REQUEST_LINK_KEY                 = byte(0x14)
	EZSP_ZIGBEE_KEY_ESTABLISHMENT_HANDLER = byte(0x9B)

	// Trust Center Frames
	EZSP_TRUST_CENTER_JOIN_HANDLER    = byte(0x24)
	EZSP_BROADCAST_NEXT_NETWORK_KEY   = byte(0x73)
	EZSP_BROADCAST_NETWORK_KEY_SWITCH = byte(0x74)
	EZSP_BECOME_TRUST_CENTER          = byte(0x77)
	EZSP_AES_MMO_HASH                 = byte(0x6F)
	EZSP_REMOVE_DEVICE                = byte(0xA8)
	EZSP_UNICAST_NWK_KEY_UPDATE       = byte(0xA9)

	// Certificate Based Key Exchange (CBKE(
	EZSP_GENERATE_CBKE_KEYS                             = byte(0xA4)
	EZSP_GENERATE_CBKE_KEYS_HANDLER                     = byte(0x9E)
	EZSP_CALCULATE_SMACS                                = byte(0x9F)
	EZSP_CALCULATE_SMACS_HANDLER                        = byte(0xA0)
	EZSP_GENERATE_CBKE_KEYS283K1                        = byte(0xE8)
	EZSP_GENERATE_CBKE_KEYS_HANDLER283K1                = byte(0xE9)
	EZSP_CALCULATE_SMACS283K1                           = byte(0xEA)
	EZSP_CALCULATE_SMACS_HANDLER283K1                   = byte(0xEB)
	EZSP_CLEAR_TEMPORARY_DATA_MAYBE_STORE_LINK_KEY      = byte(0xA1)
	EZSP_CLEAR_TEMPORARY_DATA_MAYBE_STORE_LINK_KEY283K1 = byte(0xEE)
	EZSP_GET_CERTIFICATE                                = byte(0xA5)
	EZSP_GET_CERTIFICATE283K1                           = byte(0xEC)
	EZSP_DSA_SIGN                                       = byte(0xA6)
	EZSP_DSA_SIGN_HANDLER                               = byte(0xA7)
	EZSP_DSA_VERIFY                                     = byte(0xA3)
	EZSP_DSA_VERIFY_HANDLER                             = byte(0x78)
	EZSP_SET_PREINSTALLED_CBKE_DATA                     = byte(0xA2)
	EZSP_SAVE_PREINSTALLED_CBKE_DATA283K1               = byte(0xED)

	// Mfglib
	EZSP_MFGLIB_START        = byte(0x83)
	EZSP_MFGLIB_END          = byte(0x84)
	EZSP_MFGLIB_START_TONE   = byte(0x85)
	EZSP_MFGLIB_STOP_TONE    = byte(0x86)
	EZSP_MFGLIB_START_STREAM = byte(0x87)
	EZSP_MFGLIB_STOP_STREAM  = byte(0x88)
	EZSP_MFGLIB_SEND_PACKET  = byte(0x89)
	EZSP_MFGLIB_SET_CHANNEL  = byte(0x8a)
	EZSP_MFGLIB_GET_CHANNEL  = byte(0x8b)
	EZSP_MFGLIB_SET_POWER    = byte(0x8c)
	EZSP_MFGLIB_GET_POWER    = byte(0x8d)
	EZSP_MFGLIB_RX_HANDLER   = byte(0x8e)

	// Bootloader
	EZSP_LAUNCH_STANDALONE_BOOTLOADER                     = byte(0x8f)
	EZSP_SEND_BOOTLOAD_MESSAGE                            = byte(0x90)
	EZSP_GET_STANDALONE_BOOTLOADER_VERSION_PLAT_MICRO_PHY = byte(0x91)
	EZSP_INCOMING_BOOTLOAD_MESSAGE_HANDLER                = byte(0x92)
	EZSP_BOOTLOAD_TRANSMIT_COMPLETE_HANDLER               = byte(0x93)
	EZSP_AES_ENCRYPT                                      = byte(0x94)
	EZSP_OVERRIDE_CURRENT_CHANNEL                         = byte(0x95)

	// ZLL
	EZSP_ZLL_NETWORK_OPS                = byte(0xB2)
	EZSP_ZLL_SET_INITIAL_SECURITY_STATE = byte(0xB3)
	EZSP_ZLL_START_SCAN                 = byte(0xB4)
	EZSP_ZLL_SET_RX_ON_WHEN_IDLE        = byte(0xB5)
	EZSP_ZLL_NETWORK_FOUND_HANDLER      = byte(0xB6)
	EZSP_ZLL_SCAN_COMPLETE_HANDLER      = byte(0xB7)
	EZSP_ZLL_ADDRESS_ASSIGNMENT_HANDLER = byte(0xB8)
	EZSP_SET_LOGICAL_AND_RADIO_CHANNEL  = byte(0xB9)
	EZSP_GET_LOGICAL_CHANNEL            = byte(0xBA)
	EZSP_ZLL_TOUCH_LINK_TARGET_HANDLER  = byte(0xBB)
	EZSP_ZLL_GET_TOKENS                 = byte(0xBC)
	EZSP_ZLL_SET_DATA_TOKEN             = byte(0xBD)
	EZSP_ZLL_SET_NON_ZLL_NETWORK        = byte(0xBF)
	EZSP_IS_ZLL_NETWORK                 = byte(0xBE)

	// RF4CE
	EZSP_RF4CE_SET_PAIRING_TABLE_ENTRY                  = byte(0xD0)
	EZSP_RF4CE_GET_PAIRING_TABLE_ENTRY                  = byte(0xD1)
	EZSP_RF4CE_DELETE_PAIRING_TABLE_ENTRY               = byte(0xD2)
	EZSP_RF4CE_KEY_UPDATE                               = byte(0xD3)
	EZSP_RF4CE_SEND                                     = byte(0xD4)
	EZSP_RF4CE_INCOMING_MESSAGE_HANDLER                 = byte(0xD5)
	EZSP_RF4CE_MESSAGE_SENT_HANDLER                     = byte(0xD6)
	EZSP_RF4CE_START                                    = byte(0xD7)
	EZSP_RF4CE_STOP                                     = byte(0xD8)
	EZSP_RF4CE_DISCOVERY                                = byte(0xD9)
	EZSP_RF4CE_DISCOVERY_COMPLETE_HANDLER               = byte(0xDA)
	EZSP_RF4CE_DISCOVERY_REQUEST_HANDLER                = byte(0xDB)
	EZSP_RF4CE_DISCOVERY_RESPONSE_HANDLER               = byte(0xDC)
	EZSP_RF4CE_ENABLE_AUTO_DISCOVERY_RESPONSE           = byte(0xDD)
	EZSP_RF4CE_AUTO_DISCOVERY_RESPONSE_COMPLETE_HANDLER = byte(0xDE)
	EZSP_RF4CE_PAIR                                     = byte(0xDF)
	EZSP_RF4CE_PAIR_COMPLETE_HANDLER                    = byte(0xE0)
	EZSP_RF4CE_PAIR_REQUEST_HANDLER                     = byte(0xE1)
	EZSP_RF4CE_UNPAIR                                   = byte(0xE2)
	EZSP_RF4CE_UNPAIR_HANDLER                           = byte(0xE3)
	EZSP_RF4CE_UNPAIR_COMPLETE_HANDLER                  = byte(0xE4)
	EZSP_RF4CE_SET_POWER_SAVING_PARAMETERS              = byte(0xE5)
	EZSP_RF4CE_SET_FREQUENCY_AGILITY_PARAMETERS         = byte(0xE6)
	EZSP_RF4CE_SET_APPLICATION_INFO                     = byte(0xE7)
	EZSP_RF4CE_GET_APPLICATION_INFO                     = byte(0xEF)
	EZSP_RF4CE_GET_MAX_PAYLOAD                          = byte(0xF3)
)
