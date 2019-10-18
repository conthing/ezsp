package ezsp

import "fmt"

// **************** EmberStatus ****************
const (
	EMBER_SUCCESS                                = byte(0x00)
	EMBER_ERR_FATAL                              = byte(0x01)
	EMBER_BAD_ARGUMENT                           = byte(0x02)
	EMBER_NOT_FOUND                              = byte(0x03)
	EMBER_EEPROM_MFG_STACK_VERSION_MISMATCH      = byte(0x04)
	EMBER_INCOMPATIBLE_STATIC_MEMORY_DEFINITIONS = byte(0x05)
	EMBER_EEPROM_MFG_VERSION_MISMATCH            = byte(0x06)
	EMBER_EEPROM_STACK_VERSION_MISMATCH          = byte(0x07)
	EMBER_NO_BUFFERS                             = byte(0x18)
	EMBER_SERIAL_INVALID_BAUD_RATE               = byte(0x20)
	EMBER_SERIAL_INVALID_PORT                    = byte(0x21)
	EMBER_SERIAL_TX_OVERFLOW                     = byte(0x22)
	EMBER_SERIAL_RX_OVERFLOW                     = byte(0x23)
	EMBER_SERIAL_RX_FRAME_ERROR                  = byte(0x24)
	EMBER_SERIAL_RX_PARITY_ERROR                 = byte(0x25)
	EMBER_SERIAL_RX_EMPTY                        = byte(0x26)
	EMBER_SERIAL_RX_OVERRUN_ERROR                = byte(0x27)
	EMBER_MAC_TRANSMIT_QUEUE_FULL                = byte(0x39)
	EMBER_MAC_UNKNOWN_HEADER_TYPE                = byte(0x3A)
	EMBER_MAC_ACK_HEADER_TYPE                    = byte(0x3B)
	EMBER_MAC_SCANNING                           = byte(0x3D)
	EMBER_MAC_NO_DATA                            = byte(0x31)
	EMBER_MAC_JOINED_NETWORK                     = byte(0x32)
	EMBER_MAC_BAD_SCAN_DURATION                  = byte(0x33)
	EMBER_MAC_INCORRECT_SCAN_TYPE                = byte(0x34)
	EMBER_MAC_INVALID_CHANNEL_MASK               = byte(0x35)
	EMBER_MAC_COMMAND_TRANSMIT_FAILURE           = byte(0x36)
	EMBER_MAC_NO_ACK_RECEIVED                    = byte(0x40)
	EMBER_MAC_RADIO_NETWORK_SWITCH_FAILED        = byte(0x41)
	EMBER_MAC_INDIRECT_TIMEOUT                   = byte(0x42)
	EMBER_SIM_EEPROM_ERASE_PAGE_GREEN            = byte(0x43)
	EMBER_SIM_EEPROM_ERASE_PAGE_RED              = byte(0x44)
	EMBER_SIM_EEPROM_FULL                        = byte(0x45)
	EMBER_SIM_EEPROM_INIT_1_FAILED               = byte(0x48)
	EMBER_SIM_EEPROM_INIT_2_FAILED               = byte(0x49)
	EMBER_SIM_EEPROM_INIT_3_FAILED               = byte(0x4A)
	EMBER_SIM_EEPROM_REPAIRING                   = byte(0x4D)
	EMBER_ERR_FLASH_WRITE_INHIBITED              = byte(0x46)
	EMBER_ERR_FLASH_VERIFY_FAILED                = byte(0x47)
	EMBER_ERR_FLASH_PROG_FAIL                    = byte(0x4B)
	EMBER_ERR_FLASH_ERASE_FAIL                   = byte(0x4C)
	EMBER_ERR_BOOTLOADER_TRAP_TABLE_BAD          = byte(0x58)
	EMBER_ERR_BOOTLOADER_TRAP_UNKNOWN            = byte(0x59)
	EMBER_ERR_BOOTLOADER_NO_IMAGE                = byte(0x05A)
	EMBER_DELIVERY_FAILED                        = byte(0x66)
	EMBER_BINDING_INDEX_OUT_OF_RANGE             = byte(0x69)
	EMBER_ADDRESS_TABLE_INDEX_OUT_OF_RANGE       = byte(0x6A)
	EMBER_INVALID_BINDING_INDEX                  = byte(0x6C)
	EMBER_INVALID_CALL                           = byte(0x70)
	EMBER_COST_NOT_KNOWN                         = byte(0x71)
	EMBER_MAX_MESSAGE_LIMIT_REACHED              = byte(0x72)
	EMBER_MESSAGE_TOO_LONG                       = byte(0x74)
	EMBER_BINDING_IS_ACTIVE                      = byte(0x75)
	EMBER_ADDRESS_TABLE_ENTRY_IS_ACTIVE          = byte(0x76)
	EMBER_ADC_CONVERSION_DONE                    = byte(0x80)
	EMBER_ADC_CONVERSION_BUSY                    = byte(0x81)
	EMBER_ADC_CONVERSION_DEFERRED                = byte(0x82)
	EMBER_ADC_NO_CONVERSION_PENDING              = byte(0x84)
	EMBER_SLEEP_INTERRUPTED                      = byte(0x85)
	EMBER_PHY_TX_UNDERFLOW                       = byte(0x88)
	EMBER_PHY_TX_INCOMPLETE                      = byte(0x89)
	EMBER_PHY_INVALID_CHANNEL                    = byte(0x8A)
	EMBER_PHY_INVALID_POWER                      = byte(0x8B)
	EMBER_PHY_TX_BUSY                            = byte(0x8C)
	EMBER_PHY_TX_CCA_FAIL                        = byte(0x8D)
	EMBER_PHY_OSCILLATOR_CHECK_FAILED            = byte(0x8E)
	EMBER_PHY_ACK_RECEIVED                       = byte(0x8F)
	EMBER_NETWORK_UP                             = byte(0x90)
	EMBER_NETWORK_DOWN                           = byte(0x91)
	EMBER_JOIN_FAILED                            = byte(0x94)
	EMBER_MOVE_FAILED                            = byte(0x96)
	EMBER_CANNOT_JOIN_AS_ROUTER                  = byte(0x98)
	EMBER_NODE_ID_CHANGED                        = byte(0x99)
	EMBER_PAN_ID_CHANGED                         = byte(0x9A)
	EMBER_CHANNEL_CHANGED                        = byte(0x9B)
	EMBER_NO_BEACONS                             = byte(0xAB)
	EMBER_RECEIVED_KEY_IN_THE_CLEAR              = byte(0xAC)
	EMBER_NO_NETWORK_KEY_RECEIVED                = byte(0xAD)
	EMBER_NO_LINK_KEY_RECEIVED                   = byte(0xAE)
	EMBER_PRECONFIGURED_KEY_REQUIRED             = byte(0xAF)
	EMBER_KEY_INVALID                            = byte(0xB2)
	EMBER_INVALID_SECURITY_LEVEL                 = byte(0x95)
	EMBER_APS_ENCRYPTION_ERROR                   = byte(0xA6)
	EMBER_TRUST_CENTER_MASTER_KEY_NOT_SET        = byte(0xA7)
	EMBER_SECURITY_STATE_NOT_SET                 = byte(0xA8)
	EMBER_KEY_TABLE_INVALID_ADDRESS              = byte(0xB3)
	EMBER_SECURITY_CONFIGURATION_INVALID         = byte(0xB7)
	EMBER_TOO_SOON_FOR_SWITCH_KEY                = byte(0xB8)
	EMBER_SIGNATURE_VERIFY_FAILURE               = byte(0xB9)
	EMBER_KEY_NOT_AUTHORIZED                     = byte(0xBB)
	EMBER_SECURITY_DATA_INVALID                  = byte(0xBD)
	EMBER_NOT_JOINED                             = byte(0x93)
	EMBER_NETWORK_BUSY                           = byte(0xA1)
	EMBER_INVALID_ENDPOINT                       = byte(0xA3)
	EMBER_BINDING_HAS_CHANGED                    = byte(0xA4)
	EMBER_INSUFFICIENT_RANDOM_DATA               = byte(0xA5)
	EMBER_SOURCE_ROUTE_FAILURE                   = byte(0xA9)
	EMBER_MANY_TO_ONE_ROUTE_FAILURE              = byte(0xAA)
	EMBER_STACK_AND_HARDWARE_MISMATCH            = byte(0xB0)
	EMBER_INDEX_OUT_OF_RANGE                     = byte(0xB1)
	EMBER_TABLE_FULL                             = byte(0xB4)
	EMBER_TABLE_ENTRY_ERASED                     = byte(0xB6)
	EMBER_LIBRARY_NOT_PRESENT                    = byte(0xB5)
	EMBER_OPERATION_IN_PROGRESS                  = byte(0xBA)
	EMBER_TRUST_CENTER_EUI_HAS_CHANGED           = byte(0xBC)
	EMBER_NO_RESPONSE                            = byte(0xC0)
	EMBER_DUPLICATE_ENTRY                        = byte(0xC1)
	EMBER_NOT_PERMITTED                          = byte(0xC2)
	EMBER_DISCOVERY_TIMEOUT                      = byte(0xC3)
	EMBER_DISCOVERY_ERROR                        = byte(0xC4)
	EMBER_SECURITY_TIMEOUT                       = byte(0xC5)
	EMBER_SECURITY_FAILURE                       = byte(0xC6)
	EMBER_APPLICATION_ERROR_0                    = byte(0xF0)
	EMBER_APPLICATION_ERROR_1                    = byte(0xF1)
	EMBER_APPLICATION_ERROR_2                    = byte(0xF2)
	EMBER_APPLICATION_ERROR_3                    = byte(0xF3)
	EMBER_APPLICATION_ERROR_4                    = byte(0xF4)
	EMBER_APPLICATION_ERROR_5                    = byte(0xF5)
	EMBER_APPLICATION_ERROR_6                    = byte(0xF6)
	EMBER_APPLICATION_ERROR_7                    = byte(0xF7)
	EMBER_APPLICATION_ERROR_8                    = byte(0xF8)
	EMBER_APPLICATION_ERROR_9                    = byte(0xF9)
	EMBER_APPLICATION_ERROR_10                   = byte(0xFA)
	EMBER_APPLICATION_ERROR_11                   = byte(0xFB)
	EMBER_APPLICATION_ERROR_12                   = byte(0xFC)
	EMBER_APPLICATION_ERROR_13                   = byte(0xFD)
	EMBER_APPLICATION_ERROR_14                   = byte(0xFE)
	EMBER_APPLICATION_ERROR_15                   = byte(0xFF)
)

func emberStatusToString(emberStatus byte) string {
	name, ok := emberStatusStringMap[emberStatus]
	if !ok {
		name = fmt.Sprintf("UNKNOWN_EMBERSTATUS_%02X", emberStatus)
	}
	return name
}

var emberStatusStringMap = map[byte]string{
	EMBER_SUCCESS:                                "EMBER_SUCCESS",
	EMBER_ERR_FATAL:                              "EMBER_ERR_FATAL",
	EMBER_BAD_ARGUMENT:                           "EMBER_BAD_ARGUMENT",
	EMBER_NOT_FOUND:                              "EMBER_NOT_FOUND",
	EMBER_EEPROM_MFG_STACK_VERSION_MISMATCH:      "EMBER_EEPROM_MFG_STACK_VERSION_MISMATCH",
	EMBER_INCOMPATIBLE_STATIC_MEMORY_DEFINITIONS: "EMBER_INCOMPATIBLE_STATIC_MEMORY_DEFINITIONS",
	EMBER_EEPROM_MFG_VERSION_MISMATCH:            "EMBER_EEPROM_MFG_VERSION_MISMATCH",
	EMBER_EEPROM_STACK_VERSION_MISMATCH:          "EMBER_EEPROM_STACK_VERSION_MISMATCH",
	EMBER_NO_BUFFERS:                             "EMBER_NO_BUFFERS",
	EMBER_SERIAL_INVALID_BAUD_RATE:               "EMBER_SERIAL_INVALID_BAUD_RATE",
	EMBER_SERIAL_INVALID_PORT:                    "EMBER_SERIAL_INVALID_PORT",
	EMBER_SERIAL_TX_OVERFLOW:                     "EMBER_SERIAL_TX_OVERFLOW",
	EMBER_SERIAL_RX_OVERFLOW:                     "EMBER_SERIAL_RX_OVERFLOW",
	EMBER_SERIAL_RX_FRAME_ERROR:                  "EMBER_SERIAL_RX_FRAME_ERROR",
	EMBER_SERIAL_RX_PARITY_ERROR:                 "EMBER_SERIAL_RX_PARITY_ERROR",
	EMBER_SERIAL_RX_EMPTY:                        "EMBER_SERIAL_RX_EMPTY",
	EMBER_SERIAL_RX_OVERRUN_ERROR:                "EMBER_SERIAL_RX_OVERRUN_ERROR",
	EMBER_MAC_TRANSMIT_QUEUE_FULL:                "EMBER_MAC_TRANSMIT_QUEUE_FULL",
	EMBER_MAC_UNKNOWN_HEADER_TYPE:                "EMBER_MAC_UNKNOWN_HEADER_TYPE",
	EMBER_MAC_ACK_HEADER_TYPE:                    "EMBER_MAC_ACK_HEADER_TYPE",
	EMBER_MAC_SCANNING:                           "EMBER_MAC_SCANNING",
	EMBER_MAC_NO_DATA:                            "EMBER_MAC_NO_DATA",
	EMBER_MAC_JOINED_NETWORK:                     "EMBER_MAC_JOINED_NETWORK",
	EMBER_MAC_BAD_SCAN_DURATION:                  "EMBER_MAC_BAD_SCAN_DURATION",
	EMBER_MAC_INCORRECT_SCAN_TYPE:                "EMBER_MAC_INCORRECT_SCAN_TYPE",
	EMBER_MAC_INVALID_CHANNEL_MASK:               "EMBER_MAC_INVALID_CHANNEL_MASK",
	EMBER_MAC_COMMAND_TRANSMIT_FAILURE:           "EMBER_MAC_COMMAND_TRANSMIT_FAILURE",
	EMBER_MAC_NO_ACK_RECEIVED:                    "EMBER_MAC_NO_ACK_RECEIVED",
	EMBER_MAC_RADIO_NETWORK_SWITCH_FAILED:        "EMBER_MAC_RADIO_NETWORK_SWITCH_FAILED",
	EMBER_MAC_INDIRECT_TIMEOUT:                   "EMBER_MAC_INDIRECT_TIMEOUT",
	EMBER_SIM_EEPROM_ERASE_PAGE_GREEN:            "EMBER_SIM_EEPROM_ERASE_PAGE_GREEN",
	EMBER_SIM_EEPROM_ERASE_PAGE_RED:              "EMBER_SIM_EEPROM_ERASE_PAGE_RED",
	EMBER_SIM_EEPROM_FULL:                        "EMBER_SIM_EEPROM_FULL",
	EMBER_SIM_EEPROM_INIT_1_FAILED:               "EMBER_SIM_EEPROM_INIT_1_FAILED",
	EMBER_SIM_EEPROM_INIT_2_FAILED:               "EMBER_SIM_EEPROM_INIT_2_FAILED",
	EMBER_SIM_EEPROM_INIT_3_FAILED:               "EMBER_SIM_EEPROM_INIT_3_FAILED",
	EMBER_SIM_EEPROM_REPAIRING:                   "EMBER_SIM_EEPROM_REPAIRING",
	EMBER_ERR_FLASH_WRITE_INHIBITED:              "EMBER_ERR_FLASH_WRITE_INHIBITED",
	EMBER_ERR_FLASH_VERIFY_FAILED:                "EMBER_ERR_FLASH_VERIFY_FAILED",
	EMBER_ERR_FLASH_PROG_FAIL:                    "EMBER_ERR_FLASH_PROG_FAIL",
	EMBER_ERR_FLASH_ERASE_FAIL:                   "EMBER_ERR_FLASH_ERASE_FAIL",
	EMBER_ERR_BOOTLOADER_TRAP_TABLE_BAD:          "EMBER_ERR_BOOTLOADER_TRAP_TABLE_BAD",
	EMBER_ERR_BOOTLOADER_TRAP_UNKNOWN:            "EMBER_ERR_BOOTLOADER_TRAP_UNKNOWN",
	EMBER_ERR_BOOTLOADER_NO_IMAGE:                "EMBER_ERR_BOOTLOADER_NO_IMAGE",
	EMBER_DELIVERY_FAILED:                        "EMBER_DELIVERY_FAILED",
	EMBER_BINDING_INDEX_OUT_OF_RANGE:             "EMBER_BINDING_INDEX_OUT_OF_RANGE",
	EMBER_ADDRESS_TABLE_INDEX_OUT_OF_RANGE:       "EMBER_ADDRESS_TABLE_INDEX_OUT_OF_RANGE",
	EMBER_INVALID_BINDING_INDEX:                  "EMBER_INVALID_BINDING_INDEX",
	EMBER_INVALID_CALL:                           "EMBER_INVALID_CALL",
	EMBER_COST_NOT_KNOWN:                         "EMBER_COST_NOT_KNOWN",
	EMBER_MAX_MESSAGE_LIMIT_REACHED:              "EMBER_MAX_MESSAGE_LIMIT_REACHED",
	EMBER_MESSAGE_TOO_LONG:                       "EMBER_MESSAGE_TOO_LONG",
	EMBER_BINDING_IS_ACTIVE:                      "EMBER_BINDING_IS_ACTIVE",
	EMBER_ADDRESS_TABLE_ENTRY_IS_ACTIVE:          "EMBER_ADDRESS_TABLE_ENTRY_IS_ACTIVE",
	EMBER_ADC_CONVERSION_DONE:                    "EMBER_ADC_CONVERSION_DONE",
	EMBER_ADC_CONVERSION_BUSY:                    "EMBER_ADC_CONVERSION_BUSY",
	EMBER_ADC_CONVERSION_DEFERRED:                "EMBER_ADC_CONVERSION_DEFERRED",
	EMBER_ADC_NO_CONVERSION_PENDING:              "EMBER_ADC_NO_CONVERSION_PENDING",
	EMBER_SLEEP_INTERRUPTED:                      "EMBER_SLEEP_INTERRUPTED",
	EMBER_PHY_TX_UNDERFLOW:                       "EMBER_PHY_TX_UNDERFLOW",
	EMBER_PHY_TX_INCOMPLETE:                      "EMBER_PHY_TX_INCOMPLETE",
	EMBER_PHY_INVALID_CHANNEL:                    "EMBER_PHY_INVALID_CHANNEL",
	EMBER_PHY_INVALID_POWER:                      "EMBER_PHY_INVALID_POWER",
	EMBER_PHY_TX_BUSY:                            "EMBER_PHY_TX_BUSY",
	EMBER_PHY_TX_CCA_FAIL:                        "EMBER_PHY_TX_CCA_FAIL",
	EMBER_PHY_OSCILLATOR_CHECK_FAILED:            "EMBER_PHY_OSCILLATOR_CHECK_FAILED",
	EMBER_PHY_ACK_RECEIVED:                       "EMBER_PHY_ACK_RECEIVED",
	EMBER_NETWORK_UP:                             "EMBER_NETWORK_UP",
	EMBER_NETWORK_DOWN:                           "EMBER_NETWORK_DOWN",
	EMBER_JOIN_FAILED:                            "EMBER_JOIN_FAILED",
	EMBER_MOVE_FAILED:                            "EMBER_MOVE_FAILED",
	EMBER_CANNOT_JOIN_AS_ROUTER:                  "EMBER_CANNOT_JOIN_AS_ROUTER",
	EMBER_NODE_ID_CHANGED:                        "EMBER_NODE_ID_CHANGED",
	EMBER_PAN_ID_CHANGED:                         "EMBER_PAN_ID_CHANGED",
	EMBER_CHANNEL_CHANGED:                        "EMBER_CHANNEL_CHANGED",
	EMBER_NO_BEACONS:                             "EMBER_NO_BEACONS",
	EMBER_RECEIVED_KEY_IN_THE_CLEAR:              "EMBER_RECEIVED_KEY_IN_THE_CLEAR",
	EMBER_NO_NETWORK_KEY_RECEIVED:                "EMBER_NO_NETWORK_KEY_RECEIVED",
	EMBER_NO_LINK_KEY_RECEIVED:                   "EMBER_NO_LINK_KEY_RECEIVED",
	EMBER_PRECONFIGURED_KEY_REQUIRED:             "EMBER_PRECONFIGURED_KEY_REQUIRED",
	EMBER_KEY_INVALID:                            "EMBER_KEY_INVALID",
	EMBER_INVALID_SECURITY_LEVEL:                 "EMBER_INVALID_SECURITY_LEVEL",
	EMBER_APS_ENCRYPTION_ERROR:                   "EMBER_APS_ENCRYPTION_ERROR",
	EMBER_TRUST_CENTER_MASTER_KEY_NOT_SET:        "EMBER_TRUST_CENTER_MASTER_KEY_NOT_SET",
	EMBER_SECURITY_STATE_NOT_SET:                 "EMBER_SECURITY_STATE_NOT_SET",
	EMBER_KEY_TABLE_INVALID_ADDRESS:              "EMBER_KEY_TABLE_INVALID_ADDRESS",
	EMBER_SECURITY_CONFIGURATION_INVALID:         "EMBER_SECURITY_CONFIGURATION_INVALID",
	EMBER_TOO_SOON_FOR_SWITCH_KEY:                "EMBER_TOO_SOON_FOR_SWITCH_KEY",
	EMBER_SIGNATURE_VERIFY_FAILURE:               "EMBER_SIGNATURE_VERIFY_FAILURE",
	EMBER_KEY_NOT_AUTHORIZED:                     "EMBER_KEY_NOT_AUTHORIZED",
	EMBER_SECURITY_DATA_INVALID:                  "EMBER_SECURITY_DATA_INVALID",
	EMBER_NOT_JOINED:                             "EMBER_NOT_JOINED",
	EMBER_NETWORK_BUSY:                           "EMBER_NETWORK_BUSY",
	EMBER_INVALID_ENDPOINT:                       "EMBER_INVALID_ENDPOINT",
	EMBER_BINDING_HAS_CHANGED:                    "EMBER_BINDING_HAS_CHANGED",
	EMBER_INSUFFICIENT_RANDOM_DATA:               "EMBER_INSUFFICIENT_RANDOM_DATA",
	EMBER_SOURCE_ROUTE_FAILURE:                   "EMBER_SOURCE_ROUTE_FAILURE",
	EMBER_MANY_TO_ONE_ROUTE_FAILURE:              "EMBER_MANY_TO_ONE_ROUTE_FAILURE",
	EMBER_STACK_AND_HARDWARE_MISMATCH:            "EMBER_STACK_AND_HARDWARE_MISMATCH",
	EMBER_INDEX_OUT_OF_RANGE:                     "EMBER_INDEX_OUT_OF_RANGE",
	EMBER_TABLE_FULL:                             "EMBER_TABLE_FULL",
	EMBER_TABLE_ENTRY_ERASED:                     "EMBER_TABLE_ENTRY_ERASED",
	EMBER_LIBRARY_NOT_PRESENT:                    "EMBER_LIBRARY_NOT_PRESENT",
	EMBER_OPERATION_IN_PROGRESS:                  "EMBER_OPERATION_IN_PROGRESS",
	EMBER_TRUST_CENTER_EUI_HAS_CHANGED:           "EMBER_TRUST_CENTER_EUI_HAS_CHANGED",
	EMBER_NO_RESPONSE:                            "EMBER_NO_RESPONSE",
	EMBER_DUPLICATE_ENTRY:                        "EMBER_DUPLICATE_ENTRY",
	EMBER_NOT_PERMITTED:                          "EMBER_NOT_PERMITTED",
	EMBER_DISCOVERY_TIMEOUT:                      "EMBER_DISCOVERY_TIMEOUT",
	EMBER_DISCOVERY_ERROR:                        "EMBER_DISCOVERY_ERROR",
	EMBER_SECURITY_TIMEOUT:                       "EMBER_SECURITY_TIMEOUT",
	EMBER_SECURITY_FAILURE:                       "EMBER_SECURITY_FAILURE",
	EMBER_APPLICATION_ERROR_0:                    "EMBER_APPLICATION_ERROR_0",
	EMBER_APPLICATION_ERROR_1:                    "EMBER_APPLICATION_ERROR_1",
	EMBER_APPLICATION_ERROR_2:                    "EMBER_APPLICATION_ERROR_2",
	EMBER_APPLICATION_ERROR_3:                    "EMBER_APPLICATION_ERROR_3",
	EMBER_APPLICATION_ERROR_4:                    "EMBER_APPLICATION_ERROR_4",
	EMBER_APPLICATION_ERROR_5:                    "EMBER_APPLICATION_ERROR_5",
	EMBER_APPLICATION_ERROR_6:                    "EMBER_APPLICATION_ERROR_6",
	EMBER_APPLICATION_ERROR_7:                    "EMBER_APPLICATION_ERROR_7",
	EMBER_APPLICATION_ERROR_8:                    "EMBER_APPLICATION_ERROR_8",
	EMBER_APPLICATION_ERROR_9:                    "EMBER_APPLICATION_ERROR_9",
	EMBER_APPLICATION_ERROR_10:                   "EMBER_APPLICATION_ERROR_10",
	EMBER_APPLICATION_ERROR_11:                   "EMBER_APPLICATION_ERROR_11",
	EMBER_APPLICATION_ERROR_12:                   "EMBER_APPLICATION_ERROR_12",
	EMBER_APPLICATION_ERROR_13:                   "EMBER_APPLICATION_ERROR_13",
	EMBER_APPLICATION_ERROR_14:                   "EMBER_APPLICATION_ERROR_14",
	EMBER_APPLICATION_ERROR_15:                   "EMBER_APPLICATION_ERROR_15",
}

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

// **************** Config ID ****************
const (
	// The number of packet buffers available to the stack.
	EZSP_CONFIG_PACKET_BUFFER_COUNT = byte(0x01)
	// The maximum number of router neighbors the stack can keep track of. A
	// neighbor is a node within radio range.
	EZSP_CONFIG_NEIGHBOR_TABLE_SIZE = byte(0x02)
	// The maximum number of APS retried messages the stack can be transmitting at
	// any time.
	EZSP_CONFIG_APS_UNICAST_MESSAGE_COUNT = byte(0x03)
	// The maximum number of non-volatile bindings supported by the stack.
	EZSP_CONFIG_BINDING_TABLE_SIZE = byte(0x04)
	// The maximum number of EUI64 to network address associations that the stack
	// can maintain.
	EZSP_CONFIG_ADDRESS_TABLE_SIZE = byte(0x05)
	// The maximum number of multicast groups that the device may be a member of.
	EZSP_CONFIG_MULTICAST_TABLE_SIZE = byte(0x06)
	// The maximum number of destinations to which a node can route messages. This
	// includes both messages originating at this node and those relayed for
	// others.
	EZSP_CONFIG_ROUTE_TABLE_SIZE = byte(0x07)
	// The number of simultaneous route discoveries that a node will support.
	EZSP_CONFIG_DISCOVERY_TABLE_SIZE = byte(0x08)
	// The size of the alarm broadcast buffer.
	EZSP_CONFIG_BROADCAST_ALARM_DATA_SIZE = byte(0x09)
	// The size of the unicast alarm buffers allocated for end device children.
	EZSP_CONFIG_UNICAST_ALARM_DATA_SIZE = byte(0x0A)
	// Specifies the stack profile.
	EZSP_CONFIG_STACK_PROFILE = byte(0x0C)
	// The security level used for security at the MAC and network layers. The
	// supported values are 0 (no security) and 5 (payload is encrypted and a
	// four-byte MIC is used for authentication).
	EZSP_CONFIG_SECURITY_LEVEL = byte(0x0D)
	// The maximum number of hops for a message.
	EZSP_CONFIG_MAX_HOPS = byte(0x10)
	// The maximum number of end device children that a router will support.
	EZSP_CONFIG_MAX_END_DEVICE_CHILDREN = byte(0x11)
	// The maximum amount of time that the MAC will hold a message for indirect
	// transmission to a child.
	EZSP_CONFIG_INDIRECT_TRANSMISSION_TIMEOUT = byte(0x12)
	// The maximum amount of time that an end device child can wait between polls.
	// If no poll is heard within this timeout, then the parent removes the end
	// device from its tables.
	EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT = byte(0x13)
	// The maximum amount of time that a mobile node can wait between polls. If no
	// poll is heard within this timeout, then the parent removes the mobile node
	// from its tables.
	EZSP_CONFIG_MOBILE_NODE_POLL_TIMEOUT = byte(0x14)
	// The number of child table entries reserved for use only by mobile nodes.
	EZSP_CONFIG_RESERVED_MOBILE_CHILD_ENTRIES = byte(0x15)
	// Enables boost power mode and/or the alternate transmitter output.
	EZSP_CONFIG_TX_POWER_MODE = byte(0x17)
	// 0: Allow this node to relay messages. 1: Prevent this node from relaying
	// messages.
	EZSP_CONFIG_DISABLE_RELAY = byte(0x18)
	// The maximum number of EUI64 to network address associations that the Trust
	// Center can maintain.
	EZSP_CONFIG_TRUST_CENTER_ADDRESS_CACHE_SIZE = byte(0x19)
	// The size of the source route table.
	EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE = byte(0x1A)
	// The units used for timing out end devices on their parents.
	EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT_SHIFT = byte(0x1B)
	// The number of blocks of a fragmented message that can be sent in a single
	// window.
	EZSP_CONFIG_FRAGMENT_WINDOW_SIZE = byte(0x1C)
	// The time the stack will wait (in milliseconds) between sending blocks of a
	// fragmented message.
	EZSP_CONFIG_FRAGMENT_DELAY_MS = byte(0x1D)
	// The size of the Key Table used for storing individual link keys (if the
	// device is a Trust Center) or Application Link Keys (if the device is a
	// normal node).
	EZSP_CONFIG_KEY_TABLE_SIZE = byte(0x1E)
	// The APS ACK timeout value. The stack waits this amount of time between
	// resends of APS retried messages.
	EZSP_CONFIG_APS_ACK_TIMEOUT = byte(0x1F)
	// The duration of an active scan, in the units used by the 15.4 scan
	// parameter (((1 << duration) + 1) * 15ms). This also controls the jitter
	// used when responding to a beacon request.
	EZSP_CONFIG_ACTIVE_SCAN_DURATION = byte(0x20)
	// The time the coordinator will wait (in seconds) for a second end device
	// bind request to arrive.
	EZSP_CONFIG_END_DEVICE_BIND_TIMEOUT = byte(0x21)
	// The number of PAN id conflict reports that must be received by the network
	// manager within one minute to trigger a PAN id change.
	EZSP_CONFIG_PAN_ID_CONFLICT_REPORT_THRESHOLD = byte(0x22)
	// The timeout value in minutes for how long the Trust Center or a normal node
	// waits for the ZigBee Request Key to complete. On the Trust Center this
	// controls whether or not the device buffers the request, waiting for a
	// matching pair of ZigBee Request Key. If the value is non-zero, the Trust
	// Center buffers and waits for that amount of time. If the value is zero, the
	// Trust Center does not buffer the request and immediately responds to the
	// request. Zero is the most compliant behavior.
	EZSP_CONFIG_REQUEST_KEY_TIMEOUT = byte(0x24)
	// This value indicates the size of the runtime modifiable certificate table.
	// Normally certificates are stored in MFG tokens but this table can be used
	// to field upgrade devices with new Smart Energy certificates. This value
	// cannot be set, it can only be queried.
	EZSP_CONFIG_CERTIFICATE_TABLE_SIZE = byte(0x29)
	// This is a bitmask that controls which incoming ZDO request messages are
	// passed to the application. The bits are defined in the
	// EmberZdoConfigurationFlags enumeration. To see if the application is
	// required to send a ZDO response in reply to an incoming message, the
	// application must check the APS options bitfield within the
	// incomingMessageHandler callback to see if the
	// EMBER_APS_OPTION_ZDO_RESPONSE_REQUIRED flag is set.
	EZSP_CONFIG_APPLICATION_ZDO_FLAGS = byte(0x2A)
	// The maximum number of broadcasts during a single broadcast timeout period.
	EZSP_CONFIG_BROADCAST_TABLE_SIZE = byte(0x2B)
	// The size of the MAC filter list table.
	EZSP_CONFIG_MAC_FILTER_TABLE_SIZE = byte(0x2C)
	// The number of supported networks.
	EZSP_CONFIG_SUPPORTED_NETWORKS = byte(0x2D)
	// Whether multicasts are sent to the RxOnWhenIdle=TRUE address (0xFFFD) or
	// the sleepy broadcast address (0xFFFF). The RxOnWhenIdle=TRUE address is the
	// ZigBee compliant destination for multicasts.
	EZSP_CONFIG_SEND_MULTICASTS_TO_SLEEPY_ADDRESS = byte(0x2E)
	// ZLL group address initial configuration.
	EZSP_CONFIG_ZLL_GROUP_ADDRESSES = byte(0x2F)
	// ZLL rssi threshold initial configuration.
	EZSP_CONFIG_ZLL_RSSI_THRESHOLD = byte(0x30)
	// RF4CE pairing table size.
	EZSP_CONFIG_RF4CE_PAIRING_TABLE_SIZE = byte(0x31)
	// RF4CE pending outgoing packet table size.
	EZSP_CONFIG_RF4CE_PENDING_OUTGOING_PACKET_TABLE_SIZE = byte(0x32)
	// Toggles the mtorr flow control in the stack.
	EZSP_CONFIG_MTORR_FLOW_CONTROL = byte(0x33)
	// This is a reserved frame for testing
	EZSP_CONFIG_NETWORK_TEST_PARAMETER_1 = byte(0x34)
	// This is a reserved frame for testing
	EZSP_CONFIG_NETWORK_TEST_PARAMETER_2 = byte(0x35)
)

// ID to string
func configIDToName(id byte) string {
	name, ok := configIDNameMap[id]
	if !ok {
		name = fmt.Sprintf("UNKNOWN_CONFIGID_%02X", id)
	}
	return name
}

var configIDNameMap = map[byte]string{
	EZSP_CONFIG_PACKET_BUFFER_COUNT:                      "EZSP_CONFIG_PACKET_BUFFER_COUNT",
	EZSP_CONFIG_NEIGHBOR_TABLE_SIZE:                      "EZSP_CONFIG_NEIGHBOR_TABLE_SIZE",
	EZSP_CONFIG_APS_UNICAST_MESSAGE_COUNT:                "EZSP_CONFIG_APS_UNICAST_MESSAGE_COUNT",
	EZSP_CONFIG_BINDING_TABLE_SIZE:                       "EZSP_CONFIG_BINDING_TABLE_SIZE",
	EZSP_CONFIG_ADDRESS_TABLE_SIZE:                       "EZSP_CONFIG_ADDRESS_TABLE_SIZE",
	EZSP_CONFIG_MULTICAST_TABLE_SIZE:                     "EZSP_CONFIG_MULTICAST_TABLE_SIZE",
	EZSP_CONFIG_ROUTE_TABLE_SIZE:                         "EZSP_CONFIG_ROUTE_TABLE_SIZE",
	EZSP_CONFIG_DISCOVERY_TABLE_SIZE:                     "EZSP_CONFIG_DISCOVERY_TABLE_SIZE",
	EZSP_CONFIG_BROADCAST_ALARM_DATA_SIZE:                "EZSP_CONFIG_BROADCAST_ALARM_DATA_SIZE",
	EZSP_CONFIG_UNICAST_ALARM_DATA_SIZE:                  "EZSP_CONFIG_UNICAST_ALARM_DATA_SIZE",
	EZSP_CONFIG_STACK_PROFILE:                            "EZSP_CONFIG_STACK_PROFILE",
	EZSP_CONFIG_SECURITY_LEVEL:                           "EZSP_CONFIG_SECURITY_LEVEL",
	EZSP_CONFIG_MAX_HOPS:                                 "EZSP_CONFIG_MAX_HOPS",
	EZSP_CONFIG_MAX_END_DEVICE_CHILDREN:                  "EZSP_CONFIG_MAX_END_DEVICE_CHILDREN",
	EZSP_CONFIG_INDIRECT_TRANSMISSION_TIMEOUT:            "EZSP_CONFIG_INDIRECT_TRANSMISSION_TIMEOUT",
	EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT:                  "EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT",
	EZSP_CONFIG_MOBILE_NODE_POLL_TIMEOUT:                 "EZSP_CONFIG_MOBILE_NODE_POLL_TIMEOUT",
	EZSP_CONFIG_RESERVED_MOBILE_CHILD_ENTRIES:            "EZSP_CONFIG_RESERVED_MOBILE_CHILD_ENTRIES",
	EZSP_CONFIG_TX_POWER_MODE:                            "EZSP_CONFIG_TX_POWER_MODE",
	EZSP_CONFIG_DISABLE_RELAY:                            "EZSP_CONFIG_DISABLE_RELAY",
	EZSP_CONFIG_TRUST_CENTER_ADDRESS_CACHE_SIZE:          "EZSP_CONFIG_TRUST_CENTER_ADDRESS_CACHE_SIZE",
	EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE:                  "EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE",
	EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT_SHIFT:            "EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT_SHIFT",
	EZSP_CONFIG_FRAGMENT_WINDOW_SIZE:                     "EZSP_CONFIG_FRAGMENT_WINDOW_SIZE",
	EZSP_CONFIG_FRAGMENT_DELAY_MS:                        "EZSP_CONFIG_FRAGMENT_DELAY_MS",
	EZSP_CONFIG_KEY_TABLE_SIZE:                           "EZSP_CONFIG_KEY_TABLE_SIZE",
	EZSP_CONFIG_APS_ACK_TIMEOUT:                          "EZSP_CONFIG_APS_ACK_TIMEOUT",
	EZSP_CONFIG_ACTIVE_SCAN_DURATION:                     "EZSP_CONFIG_ACTIVE_SCAN_DURATION",
	EZSP_CONFIG_END_DEVICE_BIND_TIMEOUT:                  "EZSP_CONFIG_END_DEVICE_BIND_TIMEOUT",
	EZSP_CONFIG_PAN_ID_CONFLICT_REPORT_THRESHOLD:         "EZSP_CONFIG_PAN_ID_CONFLICT_REPORT_THRESHOLD",
	EZSP_CONFIG_REQUEST_KEY_TIMEOUT:                      "EZSP_CONFIG_REQUEST_KEY_TIMEOUT",
	EZSP_CONFIG_CERTIFICATE_TABLE_SIZE:                   "EZSP_CONFIG_CERTIFICATE_TABLE_SIZE",
	EZSP_CONFIG_APPLICATION_ZDO_FLAGS:                    "EZSP_CONFIG_APPLICATION_ZDO_FLAGS",
	EZSP_CONFIG_BROADCAST_TABLE_SIZE:                     "EZSP_CONFIG_BROADCAST_TABLE_SIZE",
	EZSP_CONFIG_MAC_FILTER_TABLE_SIZE:                    "EZSP_CONFIG_MAC_FILTER_TABLE_SIZE",
	EZSP_CONFIG_SUPPORTED_NETWORKS:                       "EZSP_CONFIG_SUPPORTED_NETWORKS",
	EZSP_CONFIG_SEND_MULTICASTS_TO_SLEEPY_ADDRESS:        "EZSP_CONFIG_SEND_MULTICASTS_TO_SLEEPY_ADDRESS",
	EZSP_CONFIG_ZLL_GROUP_ADDRESSES:                      "EZSP_CONFIG_ZLL_GROUP_ADDRESSES",
	EZSP_CONFIG_ZLL_RSSI_THRESHOLD:                       "EZSP_CONFIG_ZLL_RSSI_THRESHOLD",
	EZSP_CONFIG_RF4CE_PAIRING_TABLE_SIZE:                 "EZSP_CONFIG_RF4CE_PAIRING_TABLE_SIZE",
	EZSP_CONFIG_RF4CE_PENDING_OUTGOING_PACKET_TABLE_SIZE: "EZSP_CONFIG_RF4CE_PENDING_OUTGOING_PACKET_TABLE_SIZE",
	EZSP_CONFIG_MTORR_FLOW_CONTROL:                       "EZSP_CONFIG_MTORR_FLOW_CONTROL",
	EZSP_CONFIG_NETWORK_TEST_PARAMETER_1:                 "EZSP_CONFIG_NETWORK_TEST_PARAMETER_1",
	EZSP_CONFIG_NETWORK_TEST_PARAMETER_2:                 "EZSP_CONFIG_NETWORK_TEST_PARAMETER_2",
}

// 上电默认值
const (
	DEFAULT_EZSP_CONFIG_ACTIVE_SCAN_DURATION                     = uint16(3)
	DEFAULT_EZSP_CONFIG_ADDRESS_TABLE_SIZE                       = uint16(8)
	DEFAULT_EZSP_CONFIG_APPLICATION_ZDO_FLAGS                    = uint16(0)
	DEFAULT_EZSP_CONFIG_APS_ACK_TIMEOUT                          = uint16(1600)
	DEFAULT_EZSP_CONFIG_APS_UNICAST_MESSAGE_COUNT                = uint16(10)
	DEFAULT_EZSP_CONFIG_BINDING_TABLE_SIZE                       = uint16(0)
	DEFAULT_EZSP_CONFIG_BROADCAST_ALARM_DATA_SIZE                = uint16(0)
	DEFAULT_EZSP_CONFIG_BROADCAST_TABLE_SIZE                     = uint16(15)
	DEFAULT_EZSP_CONFIG_CERTIFICATE_TABLE_SIZE                   = uint16(0)
	DEFAULT_EZSP_CONFIG_DISABLE_RELAY                            = uint16(0)
	DEFAULT_EZSP_CONFIG_DISCOVERY_TABLE_SIZE                     = uint16(8)
	DEFAULT_EZSP_CONFIG_END_DEVICE_BIND_TIMEOUT                  = uint16(60)
	DEFAULT_EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT                  = uint16(5)
	DEFAULT_EZSP_CONFIG_END_DEVICE_POLL_TIMEOUT_SHIFT            = uint16(6)
	DEFAULT_EZSP_CONFIG_FRAGMENT_DELAY_MS                        = uint16(0)
	DEFAULT_EZSP_CONFIG_FRAGMENT_WINDOW_SIZE                     = uint16(1)
	DEFAULT_EZSP_CONFIG_INDIRECT_TRANSMISSION_TIMEOUT            = uint16(3000)
	DEFAULT_EZSP_CONFIG_KEY_TABLE_SIZE                           = uint16(0)
	DEFAULT_EZSP_CONFIG_MAC_FILTER_TABLE_SIZE                    = uint16(0)
	DEFAULT_EZSP_CONFIG_MAX_END_DEVICE_CHILDREN                  = uint16(6)
	DEFAULT_EZSP_CONFIG_MAX_HOPS                                 = uint16(30)
	DEFAULT_EZSP_CONFIG_MOBILE_NODE_POLL_TIMEOUT                 = uint16(20)
	DEFAULT_EZSP_CONFIG_MTORR_FLOW_CONTROL                       = uint16(1)
	DEFAULT_EZSP_CONFIG_MULTICAST_TABLE_SIZE                     = uint16(8)
	DEFAULT_EZSP_CONFIG_NEIGHBOR_TABLE_SIZE                      = uint16(16)
	DEFAULT_EZSP_CONFIG_NETWORK_TEST_PARAMETER_1                 = uint16(8)
	DEFAULT_EZSP_CONFIG_NETWORK_TEST_PARAMETER_2                 = uint16(8)
	DEFAULT_EZSP_CONFIG_PACKET_BUFFER_COUNT                      = uint16(64)
	DEFAULT_EZSP_CONFIG_PAN_ID_CONFLICT_REPORT_THRESHOLD         = uint16(1)
	DEFAULT_EZSP_CONFIG_REQUEST_KEY_TIMEOUT                      = uint16(0)
	DEFAULT_EZSP_CONFIG_RESERVED_MOBILE_CHILD_ENTRIES            = uint16(0)
	DEFAULT_EZSP_CONFIG_RF4CE_PAIRING_TABLE_SIZE                 = uint16(4)
	DEFAULT_EZSP_CONFIG_RF4CE_PENDING_OUTGOING_PACKET_TABLE_SIZE = uint16(8)
	DEFAULT_EZSP_CONFIG_ROUTE_TABLE_SIZE                         = uint16(16)
	DEFAULT_EZSP_CONFIG_SECURITY_LEVEL                           = uint16(5)
	DEFAULT_EZSP_CONFIG_SEND_MULTICASTS_TO_SLEEPY_ADDRESS        = uint16(0)
	DEFAULT_EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE                  = uint16(0)
	DEFAULT_EZSP_CONFIG_STACK_PROFILE                            = uint16(0)
	DEFAULT_EZSP_CONFIG_SUPPORTED_NETWORKS                       = uint16(2)
	DEFAULT_EZSP_CONFIG_TRUST_CENTER_ADDRESS_CACHE_SIZE          = uint16(0)
	DEFAULT_EZSP_CONFIG_TX_POWER_MODE                            = uint16(0)
	DEFAULT_EZSP_CONFIG_UNICAST_ALARM_DATA_SIZE                  = uint16(0)
	DEFAULT_EZSP_CONFIG_ZLL_GROUP_ADDRESSES                      = uint16(1)
	DEFAULT_EZSP_CONFIG_ZLL_RSSI_THRESHOLD                       = uint16(128)
)

// **************** Value ID ****************
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

// **************** Policy ID ****************
const (
	// Controls trust center behavior.
	EZSP_TRUST_CENTER_POLICY = byte(0x00)
	// Controls how external binding modification requests are handled.
	EZSP_BINDING_MODIFICATION_POLICY = byte(0x01)
	// Controls whether the Host supplies unicast replies.
	EZSP_UNICAST_REPLIES_POLICY = byte(0x02)
	// Controls whether pollHandler callbacks are generated.
	EZSP_POLL_HANDLER_POLICY = byte(0x03)
	// Controls whether the message contents are included in the
	// messageSentHandler callback.
	EZSP_MESSAGE_CONTENTS_IN_CALLBACK_POLICY = byte(0x04)
	// Controls whether the Trust Center will respond to Trust Center link key
	// requests.
	EZSP_TC_KEY_REQUEST_POLICY = byte(0x05)
	// Controls whether the Trust Center will respond to application link key
	// requests.
	EZSP_APP_KEY_REQUEST_POLICY = byte(0x06)
	// Controls whether ZigBee packets that appear invalid are automatically
	// dropped by the stack. A counter will be incremented when this occurs.
	EZSP_PACKET_VALIDATE_LIBRARY_POLICY = byte(0x07)
	// Controls whether the stack will process ZLL messages.
	EZSP_ZLL_POLICY = byte(0x08)
	// Controls whether the ZigBee RF4CE stack will use standard profile-dependent
	// behavior during the discovery and pairing process. The profiles supported
	// at the NCP at the moment are ZRC 1.1 and MSO. If this policy is enabled the
	// stack will use standard behavior for the profiles ZRC 1.1 and MSO while it
	// will fall back to the on/off RF4CE policies for other profiles. If this
	// policy is disabled the on/off RF4CE policies are used for all profiles.
	EZSP_RF4CE_DISCOVERY_AND_PAIRING_PROFILE_BEHAVIOR_POLICY = byte(0x09)
	// Controls whether the ZigBee RF4CE stack will respond to an incoming
	// discovery request or not.
	EZSP_RF4CE_DISCOVERY_REQUEST_POLICY = byte(0x0A)
	// Controls the behavior of the ZigBee RF4CE stack discovery process.
	EZSP_RF4CE_DISCOVERY_POLICY = byte(0x0B)
	// Controls whether the ZigBee RF4CE stack will accept or deny a pair request.
	EZSP_RF4CE_PAIR_REQUEST_POLICY = byte(0x0C)
)

// **************** Decision ID ****************
const (
	// Send the network key in the clear to all joining and rejoining devices.
	EZSP_ALLOW_JOINS = byte(0x00)
	// Send the network key in the clear to all joining devices. Rejoining devices
	// are sent the network key encrypted with their trust center link key. The
	// trust center and any rejoining device are assumed to share a link key)
	// either preconfigured or obtained under a previous policy.
	EZSP_ALLOW_JOINS_REJOINS_HAVE_LINK_KEY = byte(0x04)
	// Send the network key encrypted with the joining or rejoining device's trust
	// center link key. The trust center and any joining or rejoining device are
	// assumed to share a link key, either preconfigured or obtained under a
	// previous policy. This is the default value for the
	// EZSP_TRUST_CENTER_POLICY.
	EZSP_ALLOW_PRECONFIGURED_KEY_JOINS = byte(0x01)
	// Send the network key encrypted with the rejoining device's trust center
	// link key. The trust center and any rejoining device are assumed to share a
	// link key, either preconfigured or obtained under a previous policy. No new
	// devices are allowed to join.
	EZSP_ALLOW_REJOINS_ONLY = byte(0x02)
	// Reject all unsecured join and rejoin attempts.
	EZSP_DISALLOW_ALL_JOINS_AND_REJOINS = byte(0x03)
	// EZSP_BINDING_MODIFICATION_POLICY default decision. Do not allow the local
	// binding table to be changed by remote nodes.
	EZSP_DISALLOW_BINDING_MODIFICATION = byte(0x10)
	// EZSP_BINDING_MODIFICATION_POLICY decision. Allow remote nodes to change the
	// local binding table.
	EZSP_ALLOW_BINDING_MODIFICATION = byte(0x11)
	// EZSP_BINDING_MODIFICATION_POLICY decision. Allows remote nodes to set local
	// binding entries only if the entries correspond to endpoints defined on the
	// device, and for output clusters bound to those endpoints.
	EZSP_CHECK_BINDING_MODIFICATIONS_ARE_VALID_ENDPOINT_CLUSTERS = byte(0x12)
	// EZSP_UNICAST_REPLIES_POLICY default decision. The NCP will automatically
	// send an empty reply (containing no payload) for every unicast received.
	EZSP_HOST_WILL_NOT_SUPPLY_REPLY = byte(0x20)
	// EZSP_UNICAST_REPLIES_POLICY decision. The NCP will only send a reply if it
	// receives a sendReply command from the Host.
	EZSP_HOST_WILL_SUPPLY_REPLY = byte(0x21)
	// EZSP_POLL_HANDLER_POLICY default decision. Do not inform the Host when a
	// child polls.
	EZSP_POLL_HANDLER_IGNORE = byte(0x30)
	// EZSP_POLL_HANDLER_POLICY decision. Generate a pollHandler callback when a
	// child polls.
	EZSP_POLL_HANDLER_CALLBACK = byte(0x31)
	// EZSP_MESSAGE_CONTENTS_IN_CALLBACK_POLICY default decision. Include only the
	// message tag in the messageSentHandler callback.
	EZSP_MESSAGE_TAG_ONLY_IN_CALLBACK = byte(0x40)
	// EZSP_MESSAGE_CONTENTS_IN_CALLBACK_POLICY decision. Include both the message
	// tag and the message contents in the messageSentHandler callback.
	EZSP_MESSAGE_TAG_AND_CONTENTS_IN_CALLBACK = byte(0x41)
	// EZSP_TC_KEY_REQUEST_POLICY decision. When the Trust Center receives a
	// request for a Trust Center link key, it will be ignored.
	EZSP_DENY_TC_KEY_REQUESTS = byte(0x50)
	// EZSP_TC_KEY_REQUEST_POLICY decision. When the Trust Center receives a
	// request for a Trust Center link key, it will reply to it with the
	// corresponding key.
	EZSP_ALLOW_TC_KEY_REQUESTS = byte(0x51)
	// EZSP_APP_KEY_REQUEST_POLICY decision. When the Trust Center receives a
	// request for an application link key, it will be ignored.
	EZSP_DENY_APP_KEY_REQUESTS = byte(0x60)
	// EZSP_APP_KEY_REQUEST_POLICY decision. When the Trust Center receives a
	// request for an application link key, it will randomly generate a key and
	// send it to both partners.
	EZSP_ALLOW_APP_KEY_REQUESTS = byte(0x61)
	// Indicates that packet validate library checks are enabled on the NCP.
	EZSP_PACKET_VALIDATE_LIBRARY_CHECKS_ENABLED = byte(0x62)
	// Indicates that packet validate library checks are NOT enabled on the NCP.
	EZSP_PACKET_VALIDATE_LIBRARY_CHECKS_DISABLED = byte(0x63)
	// Indicates that the RF4CE stack during discovery and pairing will use
	// standard profile-dependent behavior for the profiles ZRC 1.1 and MSO, while
	// it will fall back to the on/off policies for any other profile.
	EZSP_RF4CE_DISCOVERY_AND_PAIRING_PROFILE_BEHAVIOR_ENABLED = byte(0x70)
	// Indicates that the RF4CE stack during discovery and pairing will always use
	// the on/off policies.
	EZSP_RF4CE_DISCOVERY_AND_PAIRING_PROFILE_BEHAVIOR_DISABLED = byte(0x71)
	// Indicates that the RF4CE stack will respond to incoming discovery requests.
	EZSP_RF4CE_DISCOVERY_REQUEST_RESPOND = byte(0x72)
	// Indicates that the RF4CE stack will ignore incoming discovery requests.
	EZSP_RF4CE_DISCOVERY_REQUEST_IGNORE = byte(0x73)
	// Indicates that the RF4CE stack will perform all the discovery trials the
	// application specified in the ezspRf4ceDiscovery() call.
	EZSP_RF4CE_DISCOVERY_MAX_DISCOVERY_TRIALS = byte(0x74)
	// Indicates that the RF4CE stack will prematurely stop the discovery process
	// if a matching discovery response is received.
	EZSP_RF4CE_DISCOVERY_STOP_ON_MATCHING_RESPONSE = byte(0x75)
	// Indicates that the RF4CE stack will accept new pairings.
	EZSP_RF4CE_PAIR_REQUEST_ACCEPT = byte(0x76)
	// Indicates that the RF4CE stack will NOT accept new pairings.
	EZSP_RF4CE_PAIR_REQUEST_DENY = byte(0x77)
)

// **************** MfgToken ID ****************
const (
	// Custom version (2 bytes).
	EZSP_MFG_CUSTOM_VERSION = byte(0x00)
	// Manufacturing string (16 bytes).
	EZSP_MFG_STRING = byte(0x01)
	// Board name (16 bytes).
	EZSP_MFG_BOARD_NAME = byte(0x02)
	// Manufacturing ID (2 bytes).
	EZSP_MFG_MANUF_ID = byte(0x03)
	// Radio configuration (2 bytes).
	EZSP_MFG_PHY_CONFIG = byte(0x04)
	// Bootload AES key (16 bytes).
	EZSP_MFG_BOOTLOAD_AES_KEY = byte(0x05)
	// ASH configuration (40 bytes).
	EZSP_MFG_ASH_CONFIG = byte(0x06)
	// EZSP storage (8 bytes).
	EZSP_MFG_EZSP_STORAGE = byte(0x07)
	// Radio calibration data (64 bytes). 4 bytes are stored for each of the 16
	// channels. This token is not stored in the Flash Information Area. It is
	// updated by the stack each time a calibration is performed.
	EZSP_STACK_CAL_DATA = byte(0x08)
	// Certificate Based Key Exchange (CBKE) data (92 bytes).
	EZSP_MFG_CBKE_DATA = byte(0x09)
	// Installation code (20 bytes).
	EZSP_MFG_INSTALLATION_CODE = byte(0x0A)
	// Radio channel filter calibration data (1 byte). This token is not stored in
	// the Flash Information Area. It is updated by the stack each time a
	// calibration is performed.
	EZSP_STACK_CAL_FILTER = byte(0x0B)
	// Custom EUI64 MAC address (8 bytes).
	EZSP_MFG_CUSTOM_EUI_64 = byte(0x0C)
)

// **************** GPIO Port PIN number ****************
const (
	PORTA_PIN0 = byte((0 << 3) | 0)
	PORTA_PIN1 = byte((0 << 3) | 1)
	PORTA_PIN2 = byte((0 << 3) | 2)
	PORTA_PIN3 = byte((0 << 3) | 3)
	PORTA_PIN4 = byte((0 << 3) | 4)
	PORTA_PIN5 = byte((0 << 3) | 5)
	PORTA_PIN6 = byte((0 << 3) | 6)
	PORTA_PIN7 = byte((0 << 3) | 7)

	PORTB_PIN0 = byte((1 << 3) | 0)
	PORTB_PIN1 = byte((1 << 3) | 1)
	PORTB_PIN2 = byte((1 << 3) | 2)
	PORTB_PIN3 = byte((1 << 3) | 3)
	PORTB_PIN4 = byte((1 << 3) | 4)
	PORTB_PIN5 = byte((1 << 3) | 5)
	PORTB_PIN6 = byte((1 << 3) | 6)
	PORTB_PIN7 = byte((1 << 3) | 7)

	PORTC_PIN0 = byte((2 << 3) | 0)
	PORTC_PIN1 = byte((2 << 3) | 1)
	PORTC_PIN2 = byte((2 << 3) | 2)
	PORTC_PIN3 = byte((2 << 3) | 3)
	PORTC_PIN4 = byte((2 << 3) | 4)
	PORTC_PIN5 = byte((2 << 3) | 5)
	PORTC_PIN6 = byte((2 << 3) | 6)
	PORTC_PIN7 = byte((2 << 3) | 7)

	PORTD_PIN0 = byte((3 << 3) | 0)
	PORTD_PIN1 = byte((3 << 3) | 1)
	PORTD_PIN2 = byte((3 << 3) | 2)
	PORTD_PIN3 = byte((3 << 3) | 3)
	PORTD_PIN4 = byte((3 << 3) | 4)
	PORTD_PIN5 = byte((3 << 3) | 5)
	PORTD_PIN6 = byte((3 << 3) | 6)
	PORTD_PIN7 = byte((3 << 3) | 7)

	PORTE_PIN0 = byte((4 << 3) | 0)
	PORTE_PIN1 = byte((4 << 3) | 1)
	PORTE_PIN2 = byte((4 << 3) | 2)
	PORTE_PIN3 = byte((4 << 3) | 3)
	PORTE_PIN4 = byte((4 << 3) | 4)
	PORTE_PIN5 = byte((4 << 3) | 5)
	PORTE_PIN6 = byte((4 << 3) | 6)
	PORTE_PIN7 = byte((4 << 3) | 7)

	PORTF_PIN0 = byte((5 << 3) | 0)
	PORTF_PIN1 = byte((5 << 3) | 1)
	PORTF_PIN2 = byte((5 << 3) | 2)
	PORTF_PIN3 = byte((5 << 3) | 3)
	PORTF_PIN4 = byte((5 << 3) | 4)
	PORTF_PIN5 = byte((5 << 3) | 5)
	PORTF_PIN6 = byte((5 << 3) | 6)
	PORTF_PIN7 = byte((5 << 3) | 7)
)

// **************** Node type ****************
const (
	/** Device is not joined */
	EMBER_UNKNOWN_DEVICE = byte(0)
	/** Will relay messages and can act as a parent to other nodes. */
	EMBER_COORDINATOR = byte(1)
	/** Will relay messages and can act as a parent to other nodes. */
	EMBER_ROUTER = byte(2)
	/** Communicates only with its parent and will not relay messages. */
	EMBER_END_DEVICE = byte(3)
	/** An end device whose radio can be turned off to save power.
	 *  The application must call ::emberPollForData() to receive messages.
	 */
	EMBER_SLEEPY_END_DEVICE = byte(4)
	/** A sleepy end device that can move through the network. */
	EMBER_MOBILE_END_DEVICE = byte(5)
	/** RF4CE target node. */
	EMBER_RF4CE_TARGET = byte(6)
	/** RF4CE controller node. */
	EMBER_RF4CE_CONTROLLER = byte(7)
)

// **************** Outgoing Message type ****************
const (
	/** Unicast sent directly to an EmberNodeId. */
	EMBER_OUTGOING_DIRECT = byte(0)
	/** Unicast sent using an entry in the address table. */
	EMBER_OUTGOING_VIA_ADDRESS_TABLE = byte(1)
	/** Unicast sent using an entry in the binding table. */
	EMBER_OUTGOING_VIA_BINDING = byte(2)
	/** Multicast message.  This value is passed to emberMessageSentHandler() only.
	 * It may not be passed to emberSendUnicast(). */
	EMBER_OUTGOING_MULTICAST = byte(3)
	/** Broadcast message.  This value is passed to emberMessageSentHandler() only.
	 * It may not be passed to emberSendUnicast(). */
	EMBER_OUTGOING_BROADCAST = byte(4)
)

// **************** Incoming Message type ****************
const (
	/** Unicast. */
	EMBER_INCOMING_UNICAST = byte(0)
	/** Unicast reply. */
	EMBER_INCOMING_UNICAST_REPLY = byte(1)
	/** Multicast. */
	EMBER_INCOMING_MULTICAST = byte(2)
	/** Multicast sent by the local device. */
	EMBER_INCOMING_MULTICAST_LOOPBACK = byte(3)
	/** Broadcast. */
	EMBER_INCOMING_BROADCAST = byte(4)
	/** Broadcast sent by the local device. */
	EMBER_INCOMING_BROADCAST_LOOPBACK = byte(5)
)

// ID to string
func incomingMessageTypeToString(incomingMessageType byte) string {
	name, ok := incomingMessageTypeStringMap[incomingMessageType]
	if !ok {
		name = fmt.Sprintf("UNKNOWN_MESSAGE_TYPE_%02X", incomingMessageType)
	}
	return name
}

var incomingMessageTypeStringMap = map[byte]string{
	EMBER_INCOMING_UNICAST:            "UNICAST",
	EMBER_INCOMING_UNICAST_REPLY:      "UNICAST_REPLY",
	EMBER_INCOMING_MULTICAST:          "MULTICAST",
	EMBER_INCOMING_MULTICAST_LOOPBACK: "MULTICAST_LOOPBACK",
	EMBER_INCOMING_BROADCAST:          "BROADCAST",
	EMBER_INCOMING_BROADCAST_LOOPBACK: "BROADCAST_LOOPBACK",
}

// **************** Device Update Status sent to the Trust Center ****************
const (
	EMBER_STANDARD_SECURITY_SECURED_REJOIN   = byte(0)
	EMBER_STANDARD_SECURITY_UNSECURED_JOIN   = byte(1)
	EMBER_DEVICE_LEFT                        = byte(2)
	EMBER_STANDARD_SECURITY_UNSECURED_REJOIN = byte(3)
	EMBER_HIGH_SECURITY_SECURED_REJOIN       = byte(4)
	EMBER_HIGH_SECURITY_UNSECURED_JOIN       = byte(5)
	/* 6 Reserved */
	EMBER_HIGH_SECURITY_UNSECURED_REJOIN = byte(7)
	/* 8 - 15 Reserved */
)

// **************** Join Decision made by the Trust Center ****************
const (
	/** Allow the node to join. The node has the key. */
	EMBER_USE_PRECONFIGURED_KEY = byte(0)
	/** Allow the node to join. Send the key to the node. */
	EMBER_SEND_KEY_IN_THE_CLEAR = byte(1)
	/** Deny join. */
	EMBER_DENY_JOIN = byte(2)
	/** Take no action. */
	EMBER_NO_ACTION = byte(3)
)

// **************** Initial Security Bitmask ****************
const (
	/** This enables Distributed Trust Center Mode for the device forming the
	  network. (Previously known as ::EMBER_NO_TRUST_CENTER_MODE) */
	EMBER_DISTRIBUTED_TRUST_CENTER_MODE = uint16(0x0002)
	/** This enables a Global Link Key for the Trust Center. All nodes will share
	  the same Trust Center Link Key. */
	EMBER_TRUST_CENTER_GLOBAL_LINK_KEY = uint16(0x0004)
	/** This enables devices that perform MAC Association with a pre-configured
	  Network Key to join the network.  It is only set on the Trust Center. */
	EMBER_PRECONFIGURED_NETWORK_KEY_MODE = uint16(0x0008)

	/** This denotes that the ::EmberInitialSecurityState::preconfiguredTrustCenterEui64
	  has a value in it containing the trust center EUI64.  The device will only
	  join a network and accept commands from a trust center with that EUI64.
	  Normally this bit is NOT set, and the EUI64 of the trust center is learned
	  during the join process.  When commissioning a device to join onto
	  an existing network that is using a trust center, and without sending any
	  messages, this bit must be set and the field
	  ::EmberInitialSecurityState::preconfiguredTrustCenterEui64 must be
	  populated with the appropriate EUI64.
	*/
	EMBER_HAVE_TRUST_CENTER_EUI64 = uint16(0x0040)

	/** This denotes that the ::EmberInitialSecurityState::preconfiguredKey
	  is not the actual Link Key but a Root Key known only to the Trust Center.
	  It is hashed with the IEEE Address of the destination device in order
	  to create the actual Link Key used in encryption.  This is bit is only
	  used by the Trust Center.  The joining device need not set this.
	*/
	EMBER_TRUST_CENTER_USES_HASHED_LINK_KEY = uint16(0x0084)

	/** This denotes that the ::EmberInitialSecurityState::preconfiguredKey
	 *  element has valid data that should be used to configure the initial
	 *  security state. */
	EMBER_HAVE_PRECONFIGURED_KEY = uint16(0x0100)
	/** This denotes that the ::EmberInitialSecurityState::networkKey
	 *  element has valid data that should be used to configure the initial
	 *  security state. */
	EMBER_HAVE_NETWORK_KEY = uint16(0x0200)
	/** This denotes to a joining node that it should attempt to
	 *  acquire a Trust Center Link Key during joining. This is
	 *  necessary if the device does not have a pre-configured
	 *  key, or wants to obtain a new one (since it may be using a
	 *  well-known key during joining). */
	EMBER_GET_LINK_KEY_WHEN_JOINING = uint16(0x0400)
	/** This denotes that a joining device should only accept an encrypted
	 *  network key from the Trust Center (using its pre-configured key).
	 *  A key sent in-the-clear by the Trust Center will be rejected
	 *  and the join will fail.  This option is only valid when utilizing
	 *  a pre-configured key. */
	EMBER_REQUIRE_ENCRYPTED_KEY = uint16(0x0800)
	/** This denotes whether the device should NOT reset its outgoing frame
	 *  counters (both NWK and APS) when ::emberSetInitialSecurityState() is
	 *  called.  Normally it is advised to reset the frame counter before
	 *  joining a new network.  However in cases where a device is joining
	 *  to the same network again (but not using ::emberRejoinNetwork())
	 *  it should keep the NWK and APS frame counters stored in its tokens.
	 *
	 *  NOTE: The application is allowed to dynamically change the behavior
	 *  via EMBER_EXT_NO_FRAME_COUNTER_RESET field.
	 */
	EMBER_NO_FRAME_COUNTER_RESET = uint16(0x1000)
	/** This denotes that the device should obtain its preconfigured key from
	 *  an installation code stored in the manufacturing token.  The token
	 *  contains a value that will be hashed to obtain the actual
	 *  preconfigured key.  If that token is not valid than the call
	 *  to ::emberSetInitialSecurityState() will fail. */
	EMBER_GET_PRECONFIGURED_KEY_FROM_INSTALL_CODE = uint16(0x2000)
)

// **************** Extended Security Bitmask ****************
const (
	// If this bit is set, we set the 'key token data' field in the Initial
	// Security Bitmask to 0 (No Preconfig Key token), otherwise we leave the
	// field as it is.
	EMBER_PRECONFIG_KEY_NOT_VALID = uint16(0x0001)

	// bits 1-3 are unused.

	/** This denotes whether a joiner node (router or end-device) uses a Global
	  Link Key or a Unique Link Key. */
	EMBER_JOINER_GLOBAL_LINK_KEY = uint16(0x0010)

	/** This denotes whether the device's outgoing frame counter is allowed to
	  be reset during forming or joining. If flag is set, the outgoing frame
	  counter is not allowed to be reset. If flag is not set, the frame
	  counter is allowed to be reset. */

	EMBER_EXT_NO_FRAME_COUNTER_RESET = uint16(0x0020)

	// bit 6-7 reserved for future use (stored in TOKEN).

	/** This denotes whether a router node should discard or accept network Leave
	  Commands. */
	EMBER_NWK_LEAVE_REQUEST_NOT_ALLOWED = uint16(0x0100)

	/** This denotes whether a node is running the latest stack specification or
	  is emulating the R18 specs behavior. If this flag is enabled, a router
	  node should only send encrypted Update Device messages while the TC should
	  only accept encrypted Updated Device messages.*/
	EMBER_R18_STACK_BEHAVIOR = uint16(0x0200)

	// bit 10 and 11 are stored in RAM only.
	// bit 11 is reserved for future use.
	EMBER_VERIFY_REQUESTED_LINK_KEY = uint16(0x0400)

	// bits 12-15 are unused.
)

// **************** Network scan types ****************
const (
	// An energy scan scans each channel for its RSSI value.
	EZSP_ENERGY_SCAN = byte(0x00)
	// An active scan scans each channel for available networks.
	EZSP_ACTIVE_SCAN = byte(0x01)
)

// **************** APS option to use when sending a message ****************
const (
	/** No options. */
	EMBER_APS_OPTION_NONE = uint16(0x0000)

	EMBER_APS_OPTION_ENCRYPT_WITH_TRANSIENT_KEY = uint16(0x0001)

	/** This signs the application layer message body (APS Frame not included)
	  and appends the ECDSA signature to the end of the message.  Needed by
	  Smart Energy applications.  This requires the CBKE and ECC libraries.
	  The ::emberDsaSignHandler() function is called after DSA signing
	  is complete but before the message has been sent by the APS layer.
	  Note that when passing a buffer to the stack for DSA signing, the final
	  byte in the buffer has special significance as an indicator of how many
	  leading bytes should be ignored for signature purposes.  Refer to API
	  documentation of emberDsaSign() or the dsaSign EZSP command for further
	  details about this requirement.
	*/
	EMBER_APS_OPTION_DSA_SIGN = uint16(0x0010)
	/** Send the message using APS Encryption, using the Link Key shared
	  with the destination node to encrypt the data at the APS Level. */
	EMBER_APS_OPTION_ENCRYPTION = uint16(0x0020)
	/** Resend the message using the APS retry mechanism.  In the mesh stack,
	  this option and the enable route discovery option must be enabled for
	  an existing route to be repaired automatically. */
	EMBER_APS_OPTION_RETRY = uint16(0x0040)
	/** Send the message with the NWK 'enable route discovery' flag, which
	  causes a route discovery to be initiated if no route to the destination
	  is known.  Note that in the mesh stack, this option and the APS retry
	  option must be enabled an existing route to be repaired
	  automatically. */
	EMBER_APS_OPTION_ENABLE_ROUTE_DISCOVERY = uint16(0x0100)
	/** Send the message with the NWK 'force route discovery' flag, which causes
	  a route discovery to be initiated even if one is known. */
	EMBER_APS_OPTION_FORCE_ROUTE_DISCOVERY = uint16(0x0200)
	/** Include the source EUI64 in the network frame. */
	EMBER_APS_OPTION_SOURCE_EUI64 = uint16(0x0400)
	/** Include the destination EUI64 in the network frame. */
	EMBER_APS_OPTION_DESTINATION_EUI64 = uint16(0x0800)
	/** Send a ZDO request to discover the node ID of the destination, if it is
	  not already know. */
	EMBER_APS_OPTION_ENABLE_ADDRESS_DISCOVERY = uint16(0x1000)
	/** This message is being sent in response to a call to
	  ::emberPollHandler().  It causes the message to be sent
	  immediately instead of being queued up until the next poll from the
	  (end device) destination. */
	EMBER_APS_OPTION_POLL_RESPONSE = uint16(0x2000)
	/** This incoming message is a valid ZDO request and the application
	* is responsible for sending a ZDO response. This flag is used only
	* within emberIncomingMessageHandler() when
	* EMBER_APPLICATION_RECEIVES_UNSUPPORTED_ZDO_REQUESTS is defined. */
	EMBER_APS_OPTION_ZDO_RESPONSE_REQUIRED = uint16(0x4000)
	/** This message is part of a fragmented message.  This option may only
	  be set for unicasts.  The groupId field gives the index of this
	  fragment in the low-order byte.  If the low-order byte is zero this
	  is the first fragment and the high-order byte contains the number
	  of fragments in the message. */
	EMBER_APS_OPTION_FRAGMENT = uint16(0x8000)
)

// **************** Other const ****************
const (
	EZSP_PROTOCOL_VERSION = byte(0x04)
	EZSP_STACK_TYPE_MESH  = byte(0x02)

	EMBER_NULL_NODE_ID = uint16(0xffff)

	/** Broadcast to all routers. */
	EMBER_BROADCAST_ADDRESS = uint16(0xFFFC)
	/** Broadcast to all non-sleepy devices. */
	EMBER_RX_ON_WHEN_IDLE_BROADCAST_ADDRESS = uint16(0xFFFD)
	/** Broadcast to all devices, including sleepy end devices. */
	EMBER_SLEEPY_BROADCAST_ADDRESS = uint16(0xFFFF)

	/** @} END Broadcast Addresses */

	// From table 3.51 of 053474r14
	EMBER_MIN_BROADCAST_ADDRESS = uint16(0xFFF8)

	/**
	 * Ember Concentrator Types
	 */
	/** A concentrator with insufficient memory to store source routes for
	 * the entire network. Route records are sent to the concentrator prior
	 * to every inbound APS unicast. */
	EMBER_LOW_RAM_CONCENTRATOR = uint16(0xFFF8)
	/** A concentrator with sufficient memory to store source routes for
	 * the entire network. Remote nodes stop sending route records once
	 * the concentrator has successfully received one.
	 */
	EMBER_HIGH_RAM_CONCENTRATOR = uint16(0xFFF9)

	/**
	 * @brief Bitmask to scan recommended 802.15.4 channels.
	 */
	EMBER_RECOMMENDED_802_15_4_CHANNELS_MASK = uint32(0x0318C800)
	/**
	 * @brief Bitmask to scan all 802.15.4 channels.
	 */
	EMBER_ALL_802_15_4_CHANNELS_MASK = uint32(0x07FFF800)
	/**
	 * @brief The maximum 802.15.4 channel number is 26.
	 */
	EMBER_MAX_802_15_4_CHANNEL_NUMBER = 26

	/**
	 * @brief The minimum 802.15.4 channel number is 11.
	 */
	EMBER_MIN_802_15_4_CHANNEL_NUMBER = 11

	/**
	 * @brief There are sixteen 802.15.4 channels.
	 */
	EMBER_NUM_802_15_4_CHANNELS = EMBER_MAX_802_15_4_CHANNEL_NUMBER - EMBER_MIN_802_15_4_CHANNEL_NUMBER + 1
)
