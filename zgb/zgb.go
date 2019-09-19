package zgb

import (
	"github.com/conthing/ezsp/ash"
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

	err = ezsp.NcpConfig()
	if err != nil {
		common.Log.Errorf("NcpConfig failed: %v", err)
	}

	ezsp.NcpPrintAllConfigurations()
	//common.Log.Infof("NcpPrintAllConfigurations OK")

	eui64, err := ezsp.EzspGetEUI64()
	if err != nil {
		common.Log.Errorf("EzspGetEUI64 failed: %v", err)
	}
	common.Log.Infof("EUI64 = %016x", eui64)

	

	for {
		ezsp.NcpTick()

	}
}

//func getModuleInfo() (*models.StModuleInfo, error) {
//}
