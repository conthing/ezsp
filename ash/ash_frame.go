package ash

import (
	"fmt"

	"github.com/conthing/utils/common"
	"github.com/conthing/utils/crc16"
)

const (
	ASH_XON = byte(0x11) /*!< XON: Resume transmission
	Used in XON/XOFF flow control. Always ignored if received by the NCP. */
	ASH_XOFF = byte(0x13) /*!< XOFF: Stop transmission
	Used in XON/XOFF flow control. Always ignored if received by the NCP. */
	ASH_SUB = byte(0x18) /*!<  Substitute Byte: Replaces a byte received with a low-level communication error (e.g., framing error) from the UART.
	When a Substitute Byte is processed, the data between the previous and the next Flag Bytes is ignored. */
	ASH_CAN = byte(0x1A) /*!<  Cancel Byte: Terminates a frame in progress.
	A Cancel Byte causes all data received since the previous Flag Byte to be ignored. Note that as a special case,
	RST and RSTACK frames are preceded by Cancel Bytes to ignore any link startup noise. */
	ASH_FLAG = byte(0x7E) /*!<  Flag Byte: Marks the end of a frame.
	When a Flag Byte is received, the data received since the last Flag Byte or Cancel Byte is tested to see whether it
	is a valid frame. */
	ASH_ESC = byte(0x7D) /*!<  Escape Byte: Indicates that the following byte is escaped.
	If the byte after the Escape Byte is not a reserved byte, bit 5 of the byte is complemented to restore its original
	value. If the byte after the Escape Byte is a reserved value, the Escape Byte has no effect. */

	/*ESC后的字节与FLIP异或*/
	ASH_FLIP = byte(0x20) /*!< XOR mask used in byte stuffing */
)

var readStatusEsc = false
var readStatusSubstitute = false
var readBufferOffset = byte(0)
var readBuffer = make([]byte, 256)

func ashFrameTrace(format string, v ...interface{}) {
	//common.Log.Debugf(format, v...)
}

func ashFrameRxByteParse(recvChar byte) (err error) {
	msgDone := false

	if readStatusSubstitute {
		// ASH_SUB. ignore until next ASH_FLAG
		if recvChar == ASH_FLAG {
			readStatusSubstitute = false
		}
	} else if readStatusEsc {
		readBuffer[readBufferOffset] = (byte)(recvChar ^ ASH_FLIP)
		readBufferOffset++
		readStatusEsc = false
	} else if recvChar == ASH_ESC {
		readStatusEsc = true
	} else if recvChar == ASH_XON {
		ashFrameTrace("rx XON after: %x", readBuffer[:readBufferOffset])
		msgDone = true
	} else if recvChar == ASH_XOFF {
		common.Log.Warnf("rx XOFF after: %x", readBuffer[:readBufferOffset])
		msgDone = true
	} else if recvChar == ASH_SUB {
		common.Log.Warnf("rx SUB after: %x", readBuffer[:readBufferOffset])
		msgDone = true
		readStatusSubstitute = true
	} else if recvChar == ASH_CAN {
		ashFrameTrace("rx CANCEL after: %x", readBuffer[:readBufferOffset])
		msgDone = true
	} else if recvChar == ASH_FLAG {
		msgDone = true
		readStatusSubstitute = false

		crc16 := crc16.CRC16CCITTFalse(readBuffer[:readBufferOffset])

		if readBufferOffset <= 2 {
			common.Log.Warnf("rx frame too short < %x", readBuffer[:readBufferOffset])
			err = fmt.Errorf("rx frame too short < %x", readBuffer[:readBufferOffset])
		} else if crc16 != 0 {
			common.Log.Warnf("rx frame crc error < %x", readBuffer[:readBufferOffset])
			_ = ashRecvFrame(nil) //crc不对发送NAK
			err = fmt.Errorf("rx frame crc error < %x", readBuffer[:readBufferOffset])
		} else {
			ashFrameTrace("rx < %x", readBuffer[:readBufferOffset])
			err = ashRecvFrame(readBuffer[:readBufferOffset-2])
		}
	} else {
		readBuffer[readBufferOffset] = recvChar
		readBufferOffset++
	}
	if msgDone {
		readStatusEsc = false
		readBufferOffset = 0
	}
	return
}

func ashSendCancelByte() error {
	if ashSerial == nil {
		return fmt.Errorf("tx CANCEL failed. serial port not open")
	}

	_, err := ashSerial.Write([]byte{ASH_CAN})
	if err != nil {
		return fmt.Errorf("tx CANCEL failed. %v", err)
	}
	ashFrameTrace("tx CANCEL")
	return nil
}

func ashSendFrame(frame []byte) error {
	if ashSerial == nil {
		return fmt.Errorf("tx failed. serial port not open")
	}

	var writeBuffer []byte
	crc16 := crc16.CRC16CCITTFalse(frame)
	for _, b := range frame {
		if b == ASH_XON || b == ASH_XOFF || b == ASH_SUB || b == ASH_CAN || b == ASH_ESC || b == ASH_FLAG {
			writeBuffer = append(writeBuffer, ASH_ESC, b^ASH_FLIP)
		} else {
			writeBuffer = append(writeBuffer, b)
		}
	}
	writeBuffer = append(writeBuffer, byte((crc16>>8)&0xff), byte(crc16&0xff), ASH_FLAG)

	_, err := ashSerial.Write(writeBuffer)
	if err != nil {
		return fmt.Errorf("tx %x failed. %v", writeBuffer, err)
	}
	ashFrameTrace("tx > %x", writeBuffer)
	return nil
}
