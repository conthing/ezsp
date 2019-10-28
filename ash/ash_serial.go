package ash

import (
	"errors"
	"fmt"
	"io"

	"github.com/conthing/utils/common"
	"github.com/jacobsa/go-serial/serial"
)

var ashSerial io.ReadWriteCloser
var ashSerialXonXoff bool
var errAshRecvHandleBusy = errors.New("Recv handle busy")

//AshSerialOpen 打开串口
func AshSerialOpen(name string, baud uint, rtsCts bool) (err error) {
	options := serial.OpenOptions{
		PortName:              name,
		BaudRate:              baud,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            serial.PARITY_NONE,
		RTSCTSFlowControl:     rtsCts,
		MinimumReadSize:       0,
		InterCharacterTimeout: 50,
	}

	ashSerialXonXoff = !rtsCts

	ashSerial, err = serial.Open(options)
	if err != nil {
		ashSerial = nil
		return fmt.Errorf("failed to open serial. %v", err)
	}
	return nil
}

// AshSerialClose 关闭串口
func AshSerialClose() {
	if ashSerial != nil {
		ashSerial.Close()
	}
}

func AshSerialFlush() {
	if ashSerial != nil {
		dummy := make([]byte, 1200)
		ashSerial.Read(dummy)
	}
}

var rcvBuff = make([]byte, 1200)
var rcvStartPtr = 0

// AshSerialRecv 串口接收
func AshSerialRecv() error {
	if ashSerial == nil {
		return fmt.Errorf("failed to recv. serial port not open")
	}
	n, err := ashSerial.Read(rcvBuff[rcvStartPtr:]) //保留上次busy后剩余字节
	if n != 0 {
		len := n + rcvStartPtr
		busy := false
		offset := 0
		for i, d := range rcvBuff[:len] {
			if busy {
				rcvBuff[i-offset] = d //将后面的剩余字节搬到前面来
			} else {
				parseErr := ashFrameRxByteParse(d)
				if parseErr == errAshRecvHandleBusy {
					busy = true
					offset = i + 1
					rcvStartPtr = len - offset
					common.Log.Warnf("recv %d bytes but %d bytes remain unhandled", len, rcvStartPtr)
				} else if parseErr != nil {
					common.Log.Errorf("serial recv 0x%02x parse error %v", d, parseErr)
				}
			}

		}
	} else if err == io.EOF {
		return err
	} else if err != nil {
		return fmt.Errorf("failed to recv: %v", err)
	}
	return nil
}
