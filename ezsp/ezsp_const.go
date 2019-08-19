package ezsp

import "fmt"

// **************** EzspStatus ****************
const (
	// Success.
	EZSP_SUCCESS = byte(0x00)
	// Fatal error.
	EZSP_SPI_ERR_FATAL = byte(0x10)
	// The Response frame of the current transaction indicates the NCP has reset.
	EZSP_SPI_ERR_NCP_RESET = byte(0x11)
	// The NCP is reporting that the Command frame of the current transaction is
	// oversized (the length byte is too large).
	EZSP_SPI_ERR_OVERSIZED_EZSP_FRAME = byte(0x12)
	// The Response frame of the current transaction indicates the previous
	// transaction was aborted (nSSEL deasserted too soon).
	EZSP_SPI_ERR_ABORTED_TRANSACTION = byte(0x13)
	// The Response frame of the current transaction indicates the frame
	// terminator is missing from the Command frame.
	EZSP_SPI_ERR_MISSING_FRAME_TERMINATOR = byte(0x14)
	// The NCP has not provided a Response within the time limit defined by
	// WAIT_SECTION_TIMEOUT.
	EZSP_SPI_ERR_WAIT_SECTION_TIMEOUT = byte(0x15)
	// The Response frame from the NCP is missing the frame terminator.
	EZSP_SPI_ERR_NO_FRAME_TERMINATOR = byte(0x16)
	// The Host attempted to send an oversized Command (the length byte is too
	// large) and the AVR's spi-protocol.c blocked the transmission.
	EZSP_SPI_ERR_EZSP_COMMAND_OVERSIZED = byte(0x17)
	// The NCP attempted to send an oversized Response (the length byte is too
	// large) and the AVR's spi-protocol.c blocked the reception.
	EZSP_SPI_ERR_EZSP_RESPONSE_OVERSIZED = byte(0x18)
	// The Host has sent the Command and is still waiting for the NCP to send a
	// Response.
	EZSP_SPI_WAITING_FOR_RESPONSE = byte(0x19)
	// The NCP has not asserted nHOST_INT within the time limit defined by
	// WAKE_HANDSHAKE_TIMEOUT.
	EZSP_SPI_ERR_HANDSHAKE_TIMEOUT = byte(0x1A)
	// The NCP has not asserted nHOST_INT after an NCP reset within the time limit
	// defined by STARTUP_TIMEOUT.
	EZSP_SPI_ERR_STARTUP_TIMEOUT = byte(0x1B)
	// The Host attempted to verify the SPI Protocol activity and version number)
	// and the verification failed.
	EZSP_SPI_ERR_STARTUP_FAIL = byte(0x1C)
	// The Host has sent a command with a SPI Byte that is unsupported by the
	// current mode the NCP is operating in.
	EZSP_SPI_ERR_UNSUPPORTED_SPI_COMMAND = byte(0x1D)
	// Operation not yet complete.
	EZSP_ASH_IN_PROGRESS = byte(0x20)
	// Fatal error detected by host.
	EZSP_ASH_HOST_FATAL_ERROR = byte(0x21)
	// Fatal error detected by NCP.
	EZSP_ASH_NCP_FATAL_ERROR = byte(0x22)
	// Tried to send DATA frame too long.
	EZSP_ASH_DATA_FRAME_TOO_LONG = byte(0x23)
	// Tried to send DATA frame too short.
	EZSP_ASH_DATA_FRAME_TOO_SHORT = byte(0x24)
	// No space for tx'ed DATA frame.
	EZSP_ASH_NO_TX_SPACE = byte(0x25)
	// No space for rec'd DATA frame.
	EZSP_ASH_NO_RX_SPACE = byte(0x26)
	// No receive data available.
	EZSP_ASH_NO_RX_DATA = byte(0x27)
	// Not in Connected state.
	EZSP_ASH_NOT_CONNECTED = byte(0x28)
	// The NCP received a command before the EZSP version had been set.
	EZSP_ERROR_VERSION_NOT_SET = byte(0x30)
	// The NCP received a command containing an unsupported frame ID.
	EZSP_ERROR_INVALID_FRAME_ID = byte(0x31)
	// The direction flag in the frame control field was incorrect.
	EZSP_ERROR_WRONG_DIRECTION = byte(0x32)
	// The truncated flag in the frame control field was set, indicating there was
	// not enough memory available to complete the response or that the response
	// would have exceeded the maximum EZSP frame length.
	EZSP_ERROR_TRUNCATED = byte(0x33)
	// The overflow flag in the frame control field was set, indicating one or
	// more callbacks occurred since the previous response and there was not
	// enough memory available to report them to the Host.
	EZSP_ERROR_OVERFLOW = byte(0x34)
	// Insufficient memory was available.
	EZSP_ERROR_OUT_OF_MEMORY = byte(0x35)
	// The value was out of bounds.
	EZSP_ERROR_INVALID_VALUE = byte(0x36)
	// The configuration id was not recognized.
	EZSP_ERROR_INVALID_ID = byte(0x37)
	// Configuration values can no longer be modified.
	EZSP_ERROR_INVALID_CALL = byte(0x38)
	// The NCP failed to respond to a command.
	EZSP_ERROR_NO_RESPONSE = byte(0x39)
	// The length of the command exceeded the maximum EZSP frame length.
	EZSP_ERROR_COMMAND_TOO_LONG = byte(0x40)
	// The UART receive queue was full causing a callback response to be dropped.
	EZSP_ERROR_QUEUE_FULL = byte(0x41)
	// The command has been filtered out by NCP.
	EZSP_ERROR_COMMAND_FILTERED = byte(0x42)
	// Incompatible ASH version
	EZSP_ASH_ERROR_VERSION = byte(0x50)
	// Exceeded max ACK timeouts
	EZSP_ASH_ERROR_TIMEOUTS = byte(0x51)
	// Timed out waiting for RSTACK
	EZSP_ASH_ERROR_RESET_FAIL = byte(0x52)
	// Unexpected ncp reset
	EZSP_ASH_ERROR_NCP_RESET = byte(0x53)
	// Serial port initialization failed
	EZSP_ASH_ERROR_SERIAL_INIT = byte(0x54)
	// Invalid ncp processor type
	EZSP_ASH_ERROR_NCP_TYPE = byte(0x55)
	// Invalid ncp reset method
	EZSP_ASH_ERROR_RESET_METHOD = byte(0x56)
	// XON/XOFF not supported by host driver
	EZSP_ASH_ERROR_XON_XOFF = byte(0x57)
	// ASH protocol started
	EZSP_ASH_STARTED = byte(0x70)
	// ASH protocol connected
	EZSP_ASH_CONNECTED = byte(0x71)
	// ASH protocol disconnected
	EZSP_ASH_DISCONNECTED = byte(0x72)
	// Timer expired waiting for ack
	EZSP_ASH_ACK_TIMEOUT = byte(0x73)
	// Frame in progress cancelled
	EZSP_ASH_CANCELLED = byte(0x74)
	// Received frame out of sequence
	EZSP_ASH_OUT_OF_SEQUENCE = byte(0x75)
	// Received frame with CRC error
	EZSP_ASH_BAD_CRC = byte(0x76)
	// Received frame with comm error
	EZSP_ASH_COMM_ERROR = byte(0x77)
	// Received frame with bad ackNum
	EZSP_ASH_BAD_ACKNUM = byte(0x78)
	// Received frame shorter than minimum
	EZSP_ASH_TOO_SHORT = byte(0x79)
	// Received frame longer than maximum
	EZSP_ASH_TOO_LONG = byte(0x7A)
	// Received frame with illegal control byte
	EZSP_ASH_BAD_CONTROL = byte(0x7B)
	// Received frame with illegal length for its type
	EZSP_ASH_BAD_LENGTH = byte(0x7C)
	// No reset or error
	EZSP_ASH_NO_ERROR = byte(0xFF)
)

// ID to string
func ezspStatusToString(ezspStatus byte) string {
	name, ok := ezspStatusStringMap[ezspStatus]
	if !ok {
		name = fmt.Sprintf("UNKNOWN_EZSPSTATUS_%02X", ezspStatus)
	}
	return name
}

var ezspStatusStringMap = map[byte]string{
	EZSP_SUCCESS:                          "EZSP_SUCCESS",
	EZSP_SPI_ERR_FATAL:                    "EZSP_SPI_ERR_FATAL",
	EZSP_SPI_ERR_NCP_RESET:                "EZSP_SPI_ERR_NCP_RESET",
	EZSP_SPI_ERR_OVERSIZED_EZSP_FRAME:     "EZSP_SPI_ERR_OVERSIZED_EZSP_FRAME",
	EZSP_SPI_ERR_ABORTED_TRANSACTION:      "EZSP_SPI_ERR_ABORTED_TRANSACTION",
	EZSP_SPI_ERR_MISSING_FRAME_TERMINATOR: "EZSP_SPI_ERR_MISSING_FRAME_TERMINATOR",
	EZSP_SPI_ERR_WAIT_SECTION_TIMEOUT:     "EZSP_SPI_ERR_WAIT_SECTION_TIMEOUT",
	EZSP_SPI_ERR_NO_FRAME_TERMINATOR:      "EZSP_SPI_ERR_NO_FRAME_TERMINATOR",
	EZSP_SPI_ERR_EZSP_COMMAND_OVERSIZED:   "EZSP_SPI_ERR_EZSP_COMMAND_OVERSIZED",
	EZSP_SPI_ERR_EZSP_RESPONSE_OVERSIZED:  "EZSP_SPI_ERR_EZSP_RESPONSE_OVERSIZED",
	EZSP_SPI_WAITING_FOR_RESPONSE:         "EZSP_SPI_WAITING_FOR_RESPONSE",
	EZSP_SPI_ERR_HANDSHAKE_TIMEOUT:        "EZSP_SPI_ERR_HANDSHAKE_TIMEOUT",
	EZSP_SPI_ERR_STARTUP_TIMEOUT:          "EZSP_SPI_ERR_STARTUP_TIMEOUT",
	EZSP_SPI_ERR_STARTUP_FAIL:             "EZSP_SPI_ERR_STARTUP_FAIL",
	EZSP_SPI_ERR_UNSUPPORTED_SPI_COMMAND:  "EZSP_SPI_ERR_UNSUPPORTED_SPI_COMMAND",
	EZSP_ASH_IN_PROGRESS:                  "EZSP_ASH_IN_PROGRESS",
	EZSP_ASH_HOST_FATAL_ERROR:             "EZSP_ASH_HOST_FATAL_ERROR",
	EZSP_ASH_NCP_FATAL_ERROR:              "EZSP_ASH_NCP_FATAL_ERROR",
	EZSP_ASH_DATA_FRAME_TOO_LONG:          "EZSP_ASH_DATA_FRAME_TOO_LONG",
	EZSP_ASH_DATA_FRAME_TOO_SHORT:         "EZSP_ASH_DATA_FRAME_TOO_SHORT",
	EZSP_ASH_NO_TX_SPACE:                  "EZSP_ASH_NO_TX_SPACE",
	EZSP_ASH_NO_RX_SPACE:                  "EZSP_ASH_NO_RX_SPACE",
	EZSP_ASH_NO_RX_DATA:                   "EZSP_ASH_NO_RX_DATA",
	EZSP_ASH_NOT_CONNECTED:                "EZSP_ASH_NOT_CONNECTED",
	EZSP_ERROR_VERSION_NOT_SET:            "EZSP_ERROR_VERSION_NOT_SET",
	EZSP_ERROR_INVALID_FRAME_ID:           "EZSP_ERROR_INVALID_FRAME_ID",
	EZSP_ERROR_WRONG_DIRECTION:            "EZSP_ERROR_WRONG_DIRECTION",
	EZSP_ERROR_TRUNCATED:                  "EZSP_ERROR_TRUNCATED",
	EZSP_ERROR_OVERFLOW:                   "EZSP_ERROR_OVERFLOW",
	EZSP_ERROR_OUT_OF_MEMORY:              "EZSP_ERROR_OUT_OF_MEMORY",
	EZSP_ERROR_INVALID_VALUE:              "EZSP_ERROR_INVALID_VALUE",
	EZSP_ERROR_INVALID_ID:                 "EZSP_ERROR_INVALID_ID",
	EZSP_ERROR_INVALID_CALL:               "EZSP_ERROR_INVALID_CALL",
	EZSP_ERROR_NO_RESPONSE:                "EZSP_ERROR_NO_RESPONSE",
	EZSP_ERROR_COMMAND_TOO_LONG:           "EZSP_ERROR_COMMAND_TOO_LONG",
	EZSP_ERROR_QUEUE_FULL:                 "EZSP_ERROR_QUEUE_FULL",
	EZSP_ERROR_COMMAND_FILTERED:           "EZSP_ERROR_COMMAND_FILTERED",
	EZSP_ASH_ERROR_VERSION:                "EZSP_ASH_ERROR_VERSION",
	EZSP_ASH_ERROR_TIMEOUTS:               "EZSP_ASH_ERROR_TIMEOUTS",
	EZSP_ASH_ERROR_RESET_FAIL:             "EZSP_ASH_ERROR_RESET_FAIL",
	EZSP_ASH_ERROR_NCP_RESET:              "EZSP_ASH_ERROR_NCP_RESET",
	EZSP_ASH_ERROR_SERIAL_INIT:            "EZSP_ASH_ERROR_SERIAL_INIT",
	EZSP_ASH_ERROR_NCP_TYPE:               "EZSP_ASH_ERROR_NCP_TYPE",
	EZSP_ASH_ERROR_RESET_METHOD:           "EZSP_ASH_ERROR_RESET_METHOD",
	EZSP_ASH_ERROR_XON_XOFF:               "EZSP_ASH_ERROR_XON_XOFF",
	EZSP_ASH_STARTED:                      "EZSP_ASH_STARTED",
	EZSP_ASH_CONNECTED:                    "EZSP_ASH_CONNECTED",
	EZSP_ASH_DISCONNECTED:                 "EZSP_ASH_DISCONNECTED",
	EZSP_ASH_ACK_TIMEOUT:                  "EZSP_ASH_ACK_TIMEOUT",
	EZSP_ASH_CANCELLED:                    "EZSP_ASH_CANCELLED",
	EZSP_ASH_OUT_OF_SEQUENCE:              "EZSP_ASH_OUT_OF_SEQUENCE",
	EZSP_ASH_BAD_CRC:                      "EZSP_ASH_BAD_CRC",
	EZSP_ASH_COMM_ERROR:                   "EZSP_ASH_COMM_ERROR",
	EZSP_ASH_BAD_ACKNUM:                   "EZSP_ASH_BAD_ACKNUM",
	EZSP_ASH_TOO_SHORT:                    "EZSP_ASH_TOO_SHORT",
	EZSP_ASH_TOO_LONG:                     "EZSP_ASH_TOO_LONG",
	EZSP_ASH_BAD_CONTROL:                  "EZSP_ASH_BAD_CONTROL",
	EZSP_ASH_BAD_LENGTH:                   "EZSP_ASH_BAD_LENGTH",
	EZSP_ASH_NO_ERROR:                     "EZSP_ASH_NO_ERROR",
}

// **************** Frame ID ****************
const (
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

// ID to string
func frameIDToName(id byte) string {
	name, ok := frameIDNameMap[id]
	if !ok {
		name = fmt.Sprintf("UNKNOWN_ID_%02X", id)
	}
	return name
}

var frameIDNameMap = map[byte]string{
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

// 判断是否callback
func isValidCallbackID(callbackID byte) bool {
	if isCallbackIDMap[allCallbackIDs[0]] == false {
		for _, id := range allCallbackIDs {
			isCallbackIDMap[id] = true
		}
	}
	return isCallbackIDMap[callbackID]
}

var isCallbackIDMap [256]bool
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

// **************** EzspGetValue ID ****************
const (
	// The contents of the node data stack token.
	EZSP_VALUE_TOKEN_STACK_NODE_DATA = byte(0x00)
	// The types of MAC passthrough messages that the host wishes to receive.
	EZSP_VALUE_MAC_PASSTHROUGH_FLAGS = byte(0x01)
	// The source address used to filter legacy EmberNet messages when the
	// EMBER_MAC_PASSTHROUGH_EMBERNET_SOURCE flag is set in
	// EZSP_VALUE_MAC_PASSTHROUGH_FLAGS.
	EZSP_VALUE_EMBERNET_PASSTHROUGH_SOURCE_ADDRESS = byte(0x02)
	// The number of available message buffers.
	EZSP_VALUE_FREE_BUFFERS = byte(0x03)
	// Selects sending synchronous callbacks in ezsp-uart.
	EZSP_VALUE_UART_SYNCH_CALLBACKS = byte(0x04)
	// The maximum incoming transfer size for the local node.
	EZSP_VALUE_MAXIMUM_INCOMING_TRANSFER_SIZE = byte(0x05)
	// The maximum outgoing transfer size for the local node.
	EZSP_VALUE_MAXIMUM_OUTGOING_TRANSFER_SIZE = byte(0x06)
	// A boolean indicating whether stack tokens are written to persistent storage
	// as they change.
	EZSP_VALUE_STACK_TOKEN_WRITING = byte(0x07)
	// A read-only value indicating whether the stack is currently performing a
	// rejoin.
	EZSP_VALUE_STACK_IS_PERFORMING_REJOIN = byte(0x08)
	// A list of EmberMacFilterMatchData values.
	EZSP_VALUE_MAC_FILTER_LIST = byte(0x09)
	// The Ember Extended Security Bitmask.
	EZSP_VALUE_EXTENDED_SECURITY_BITMASK = byte(0x0A)
	// The node short ID.
	EZSP_VALUE_NODE_SHORT_ID = byte(0x0B)
	// The descriptor capability of the local node.
	EZSP_VALUE_DESCRIPTOR_CAPABILITY = byte(0x0C)
	// The stack device request sequence number of the local node.
	EZSP_VALUE_STACK_DEVICE_REQUEST_SEQUENCE_NUMBER = byte(0x0D)
	// Enable or disable radio hold-off.
	EZSP_VALUE_RADIO_HOLD_OFF = byte(0x0E)
	// The flags field associated with the endpoint data.
	EZSP_VALUE_ENDPOINT_FLAGS = byte(0x0F)
	// Enable/disable the Mfg security config key settings.
	EZSP_VALUE_MFG_SECURITY_CONFIG = byte(0x10)
	// Retrieves the version information from the stack on the NCP.
	EZSP_VALUE_VERSION_INFO = byte(0x11)
	// This will get/set the rejoin reason noted by the host for a subsequent call
	// to emberFindAndRejoinNetwork(). After a call to emberFindAndRejoinNetwork()
	// the host's rejoin reason will be set to EMBER_REJOIN_REASON_NONE. The NCP
	// will store the rejoin reason used by the call to
	// emberFindAndRejoinNetwork()
	EZSP_VALUE_NEXT_HOST_REJOIN_REASON = byte(0x12)
	// This is the reason that the last rejoin took place. This value may only be
	// retrieved, not set. The rejoin may have been initiated by the stack (NCP)
	// or the application (host). If a host initiated a rejoin the reason will be
	// set by default to EMBER_REJOIN_DUE_TO_APP_EVENT_1. If the application
	// wishes to denote its own rejoin reasons it can do so by calling
	// ezspSetValue(EMBER_VALUE_HOST_REJOIN_REASON)
	// EMBER_REJOIN_DUE_TO_APP_EVENT_X). X is a number corresponding to one of the
	// app events defined. If the NCP initiated a rejoin it will record this value
	// internally for retrieval by ezspGetValue(EZSP_VALUE_REAL_REJOIN_REASON).
	EZSP_VALUE_LAST_REJOIN_REASON = byte(0x13)
	// The next ZigBee sequence number.
	EZSP_VALUE_NEXT_ZIGBEE_SEQUENCE_NUMBER = byte(0x14)
	// CCA energy detect threshold for radio.
	EZSP_VALUE_CCA_THRESHOLD = byte(0x15)
	// The RF4CE discovery LQI threshold parameter.
	EZSP_VALUE_RF4CE_DISCOVERY_LQI_THRESHOLD = byte(0x16)
	// The threshold value for a counter
	EZSP_VALUE_SET_COUNTER_THRESHOLD = byte(0x17)
	// Resets all counters thresholds to 0xFF
	EZSP_VALUE_RESET_COUNTER_THRESHOLDS = byte(0x18)
	// Clears all the counters
	EZSP_VALUE_CLEAR_COUNTERS = byte(0x19)
	// The node's new certificate signed by the CA.
	EZSP_VALUE_CERTIFICATE_283K1 = byte(0x1A)
	// The Certificate Authority's public key.
	EZSP_VALUE_PUBLIC_KEY_283K1 = byte(0x1B)
	// The node's new static private key.
	EZSP_VALUE_PRIVATE_KEY_283K1 = byte(0x1C)
	// The GDP binding recipient parameters
	EZSP_VALUE_RF4CE_GDP_BINDING_RECIPIENT_PARAMETERS = byte(0x1D)
	// The GDP binding push button stimulus received pending flag
	EZSP_VALUE_RF4CE_GDP_PUSH_BUTTON_STIMULUS_RECEIVED_PENDING_FLAG = byte(0x1E)
	// The GDP originator proxy flag in the advanced binding options
	EZSP_VALUE_RF4CE_GDP_BINDING_PROXY_FLAG = byte(0x1F)
	// The GDP application specific user string
	EZSP_VALUE_RF4CE_GDP_APPLICATION_SPECIFIC_USER_STRING = byte(0x20)
	// The MSO user string
	EZSP_VALUE_RF4CE_MSO_USER_STRING = byte(0x21)
	// The MSO binding recipient parameters
	EZSP_VALUE_RF4CE_MSO_BINDING_RECIPIENT_PARAMETERS = byte(0x22)
	// The NWK layer security frame counter value
	EZSP_VALUE_NWK_FRAME_COUNTER = byte(0x23)
	// The APS layer security frame counter value
	EZSP_VALUE_APS_FRAME_COUNTER = byte(0x24)
	// Sets the device type to use on the next rejoin using device type
	EZSP_VALUE_RETRY_DEVICE_TYPE = byte(0x25)
	// The device RF4CE base channel
	EZSP_VALUE_RF4CE_BASE_CHANNEL = byte(0x26)
	// The RF4CE device types supported by the node
	EZSP_VALUE_RF4CE_SUPPORTED_DEVICE_TYPES_LIST = byte(0x27)
	// The RF4CE profiles supported by the node
	EZSP_VALUE_RF4CE_SUPPORTED_PROFILES_LIST = byte(0x28)
)
