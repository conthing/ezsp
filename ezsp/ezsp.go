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

	err = EzspSetPolicy(EZSP_MESSAGE_CONTENTS_IN_CALLBACK_POLICY, EZSP_MESSAGE_TAG_AND_CONTENTS_IN_CALLBACK)
	if err != nil {
		return fmt.Errorf("EzspSetPolicy failed: %v", err)
	}
	common.Log.Infof("EzspSetPolicy EZSP_MESSAGE_TAG_AND_CONTENTS_IN_CALLBACK")

	err = EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE(84)
	if err != nil {
		return fmt.Errorf("EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE failed: %v", err)
	}
	common.Log.Infof("EzspSetValue_MAXIMUM_INCOMING_TRANSFER_SIZE = 84")

	err = EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE(84)
	if err != nil {
		return fmt.Errorf("EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE failed: %v", err)
	}
	common.Log.Infof("EzspSetValue_MAXIMUM_OUTGOING_TRANSFER_SIZE = 84")

	err = ncpSetRadio()
	if err != nil {
		return fmt.Errorf("ncpSetRadio failed: %v", err)
	}
	common.Log.Infof("ncpSetRadio OK")

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

//func NcpNetworkInit() (err error) {
//	nodeType, parameters, err := EzspGetNetworkParameters()
//	if err != nil {
//		return fmt.Errorf("EzspGetNetworkParameters() failed: %v", err)
//	}
//
//	var callsetup bool
//	if nodeType != EMBER_COORDINATOR {
//		callsetup = true;
//		common.Log.Infof("nodeType(%d) != EMBER_COORDINATOR, setup node", nodeType)
//	} else {
//		err = EzspNetworkInit()
//		if err != nil {
//			callsetup = true;
//			common.Log.Errorf("EzspNetworkInit() failed: %v, setup node", err)
//		}
//	}
//}

func NcpTick() {
	select {
	case cb := <-callbackCh:
		EzspCallbackDispatch(cb)
	}
}
