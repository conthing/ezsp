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

	//common.Log.Infof("%v", stackVersion)

	common.Log.Infof("NcpGetVersion: protocolVersion(%d) stackType(%d) stackVersion(%s)", ncpProtocolVersion, ncpStackType, ncpStackVersion)
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
			common.Log.Infof("%s = %d", name, value)
		}
	}
}

type EzspConfig struct {
	configID byte
	value    uint16
}

var ncpAllConfigurations = [...]EzspConfig{
	{EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE, uint16(2)},
	{EZSP_CONFIG_SECURITY_LEVEL, uint16(5)},
	{EZSP_CONFIG_ADDRESS_TABLE_SIZE, uint16(2)},
	{EZSP_CONFIG_TRUST_CENTER_ADDRESS_CACHE_SIZE, uint16(2)},
	{EZSP_CONFIG_STACK_PROFILE, uint16(0)},
	{EZSP_CONFIG_INDIRECT_TRANSMISSION_TIMEOUT, uint16(7680)},
	{EZSP_CONFIG_MAX_HOPS, uint16(30)},
	{EZSP_CONFIG_SUPPORTED_NETWORKS, uint16(1)},
}

func NcpSetConfigurations() (err error) {
	for _, cfg := range ncpAllConfigurations {
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
		common.Log.Infof("%s = %d write success", name, cfg.value)
	}
	return
}

func NcpTick() {
	select {
	case cb := <-callbackCh:
		EzspCallbackDispatch(cb)
	}
}
