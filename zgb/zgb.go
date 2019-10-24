package zgb

import (
	"github.com/conthing/ezsp/ash"
	"github.com/conthing/ezsp/c4"
	"github.com/conthing/ezsp/ezsp"
	"github.com/conthing/utils/common"
)

//TickRunning 定时运行tick
func TickRunning(errs chan error) {
	ash.AshStartTransceiver(ezsp.AshRecvImp, errs)

	err := ash.AshReset()
	if err != nil {
		common.Log.Errorf("AshReset failed: %v", err)
	} else {
		c4.C4Init()
		ezsp.EzspFrameInitVariables() // 有些变量在ASH的接收线程里会被使用
		ash.InitVariables()           // 上层的变量初始化完成后，最后调用ASH的变量初始化
		common.Log.Info("AshReset OK")
	}

	err = ezsp.NcpGetVersion()
	if err != nil {
		common.Log.Errorf("NcpGetVersion failed: %v", err)
	}

	common.Log.Infof("NCP module info : %+v", ezsp.ModuleInfo)

	err = ezsp.NcpConfig()
	if err != nil {
		common.Log.Errorf("NcpConfig failed: %v", err)
	}

	//common.Log.Infof("Print All Configurations...")
	//ezsp.NcpPrintAllConfigurations()

	rebootCnt, err := ezsp.NcpGetAndIncRebootCnt()
	if err != nil {
		common.Log.Errorf("NcpGetAndIncRebootCnt failed: %v", err)
	}
	common.Log.Infof("NCP reboot count = %d", rebootCnt)

	eui64, err := ezsp.EzspGetEUI64()
	if err != nil {
		common.Log.Errorf("EzspGetEUI64 failed: %v", err)
	}
	common.Log.Infof("NCP EUI64 = %016x", eui64)

	//err = ezsp.NcpFormNetwork(0xff)
	//if err != nil {
	//	common.Log.Errorf("NcpFormNetwork failed: %v", err)
	//}
	//common.Log.Infof("NcpFormNetwork OK")

	err = ezsp.EzspNetworkInit()
	if err != nil {
		common.Log.Errorf("EzspNetworkInit failed: %v", err)
	}
	common.Log.Infof("EzspNetworkInit OK")

	//err = c4.SetPermission(&c4.StPermission{60, []*c4.StPassport{&c4.StPassport{PS: "inSona:IN-C01-WR-4", MAC: "xxxxxxxxxxxxce73"}}})
	//if err != nil {
	//	common.Log.Errorf("C4SetPermission failed: %v", err)
	//}
	//common.Log.Infof("C4SetPermission for 60 seconds")

	for {
		c4.C4Tick()
	}
}
