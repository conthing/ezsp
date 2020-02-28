package zgb

import (
	"github.com/conthing/ezsp/ash"
	"github.com/conthing/ezsp/c4"
	"github.com/conthing/ezsp/ezsp"

	"github.com/conthing/utils/common"
)

type StTraceSettings struct {
	AshFrameTraceOn       bool
	AshTraceOn            bool
	EzspFrameTraceOn      bool
	EzspApiTraceOn        bool
	EzspCallbackTraceOn   bool
	NcpTraceOn            bool
	NcpFormTraceOn        bool
	NcpSourceRouteTraceOn bool
}

func TraceSet(settings *StTraceSettings) {
	if settings.AshFrameTraceOn {
		common.Log.Info("AshFrameTraceOn")
		ash.AshFrameTraceOn = true
	} else {
		ash.AshFrameTraceOn = false
	}
	if settings.AshTraceOn {
		common.Log.Info("AshTraceOn")
		ash.AshTraceOn = true
	} else {
		ash.AshTraceOn = false
	}
	if settings.EzspFrameTraceOn {
		common.Log.Info("EzspFrameTraceOn")
		ezsp.EzspFrameTraceOn = true
	} else {
		ezsp.EzspFrameTraceOn = false
	}
	if settings.EzspApiTraceOn {
		common.Log.Info("EzspApiTraceOn")
		ezsp.EzspApiTraceOn = true
	} else {
		ezsp.EzspApiTraceOn = false
	}
	if settings.EzspCallbackTraceOn {
		common.Log.Info("EzspCallbackTraceOn")
		ezsp.EzspCallbackTraceOn = true
	} else {
		ezsp.EzspCallbackTraceOn = false
	}
	if settings.NcpTraceOn {
		common.Log.Info("NcpTraceOn")
		ezsp.NcpTraceOn = true
	} else {
		ezsp.NcpTraceOn = false
	}
	if settings.NcpFormTraceOn {
		common.Log.Info("NcpFormTraceOn")
		ezsp.NcpFormTraceOn = true
	} else {
		ezsp.NcpFormTraceOn = false
	}
	if settings.NcpSourceRouteTraceOn {
		common.Log.Info("NcpSourceRouteTraceOn")
		ezsp.NcpSourceRouteTraceOn = true
	} else {
		ezsp.NcpSourceRouteTraceOn = false
	}
}

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
