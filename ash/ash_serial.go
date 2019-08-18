package ash

import (
	"fmt"
	"io"

	"github.com/conthing/utils/common"
	"github.com/jacobsa/go-serial/serial"
)

var ashSerial io.ReadWriteCloser
var ashSerialXonXoff bool

//AshSerialOpen 打开串口
func AshSerialOpen(name string, baud uint, rtsCts bool) (err error) {
	options := serial.OpenOptions{
		PortName:              name,
		BaudRate:              baud,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            serial.PARITY_NONE,
		RTSCTSFlowControl:     rtsCts,
		MinimumReadSize:       1,
		InterCharacterTimeout: 200,
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

// AshSerialRecv 串口接收
func AshSerialRecv() error {
	if ashSerial == nil {
		return fmt.Errorf("failed to recv. serial port not open")
	}
	data := make([]byte, 128)
	n, err := ashSerial.Read(data)
	if n != 0 {
		for _, d := range data[:n] {
			parseErr := ashFrameRxByteParse(d)
			if parseErr != nil {
				common.Log.Errorf("serial recv 0x%02x parse error %v", d, parseErr)
			}
		}
		//} else if err == io.EOF {
		//log.Print("EOF")
		//} else if err == serial.ErrTimeout {
		//log.Print("ErrTimeout")
	} else if err != nil {
		return fmt.Errorf("failed to recv: %v", err)
	}
	return nil
}
