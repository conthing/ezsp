package ezsp

import (
	"fmt"

	"github.com/conthing/utils/common"
)

var ncpProtocolVersion byte
var ncpStackType byte
var ncpStackVersion string

func NcpGetVersion() (err error) {
	var stackVersion uint16
	ncpProtocolVersion, ncpStackType, stackVersion, err = EzspVersion(EZSP_PROTOCOL_VERSION)
	if err != nil {
		return fmt.Errorf("EzspVersion failed: %v", err)
	}

	emberVersion, err := EzspGetValue_VERSION_INFO()
	if err != nil {
		common.Log.Errorf("EzspGetValue_VERSION_INFO failed: %v", err)
		ncpStackVersion = fmt.Sprintf("%d.%d.%d.%d", (stackVersion>>12)&0xF, (stackVersion>>8)&0xF, (stackVersion>>4)&0xF, stackVersion&0xF)
	} else {
		ncpStackVersion = emberVersion.String()
	}

	common.Log.Infof("NcpGetVersion: protocolVersion(%d) stackType(%d) stackVersion(%s)", ncpProtocolVersion, ncpStackType, ncpStackVersion)
	return nil
}

func NcpPrintAllConfigurations() {
	for configId, name := range configIDNameMap {
		value, err := EzspGetConfigurationValue(configId)
		if err != nil {
			common.Log.Errorf("EZSP get %s failed: %v", name, err)
		}
		common.Log.Infof("EZSP config %s = %d", name, value)
	}

}

func NcpTick() {
	select {
	case cb := <-callbackCh:
		EzspCallbackDispatch(cb)
	}
}
