package zgb

import (
	"github.com/conthing/ezsp/ash"
	"github.com/conthing/ezsp/c4"
	"github.com/conthing/ezsp/ezsp"
	"github.com/conthing/utils/common"
)

//TickRunning 定时运行tick
func TickRunning(_ chan error) {

	err := ash.AshReset()
	if err != nil {
		common.Log.Errorf("AshReset failed: %v", err)
	} else {
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

	c4.C4Init()

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

	err = ezsp.EzspPermitJoining(60)
	if err != nil {
		common.Log.Errorf("EzspPermitJoining failed: %v", err)
	}
	common.Log.Infof("EzspPermitJoining for 60 seconds")

	for {
		ezsp.EzspTick()
	}
}

//func getModuleInfo() (*models.StModuleInfo, error) {
//}
