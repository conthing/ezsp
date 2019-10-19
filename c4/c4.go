package c4

import (
	"fmt"
	"sync"
	"time"

	ezsp "github.com/conthing/ezsp/ezsp"
	"github.com/conthing/ezsp/zcl"
	"github.com/conthing/utils/common"
)

//todo 收发通知，node的online offline deleted通知
const (
	C4_PROFILE  = uint16(0xc25d)
	C4_CLUSTER  = uint16(0x0001)
	ZDO_PROFILE = uint16(0x0000)

	C4_ATTR_DEVICE_TYPE                       = uint16(0x0000)
	C4_ATTR_SSCP_ANNOUNCE_WINDOW              = uint16(0x0001)
	C4_ATTR_SSCP_MTORR_PERIOD                 = uint16(0x0002)
	C4_ATTR_SSCP_NUM_ZAPS                     = uint16(0x0003)
	C4_ATTR_FIRMWARE_VERSION                  = uint16(0x0004)
	C4_ATTR_REFLASH_VERSION                   = uint16(0x0005)
	C4_ATTR_BOOT_COUNT                        = uint16(0x0006)
	C4_ATTR_PRODUCT_STRING                    = uint16(0x0007)
	C4_ATTR_ACCESS_POINT_NODE_ID              = uint16(0x0008)
	C4_ATTR_ACCESS_POINT_LONG_ID              = uint16(0x0009)
	C4_ATTR_ACCESS_POINT_COST                 = uint16(0x000a)
	C4_ATTR_END_NODE_ACCESS_POINT_POLL_PERIOD = uint16(0x000b)
	C4_ATTR_SSCP_CHANNEL                      = uint16(0x000c)
	C4_ATTR_ZSERVER_IP                        = uint16(0x000d)
	C4_ATTR_ZSERVER_HOST_NAME                 = uint16(0x000e)
	C4_ATTR_END_NODE_PARENT_POLL_PERIOD       = uint16(0x000f)
	C4_ATTR_MESH_LONG_ID                      = uint16(0x0010)
	C4_ATTR_MESH_SHORT_ID                     = uint16(0x0011)
	C4_ATTR_AP_TABLE                          = uint16(0x0012)
	C4_ATTR_AVG_RSSI                          = uint16(0x0013)
	C4_ATTR_AVG_LQI                           = uint16(0x0014)
	C4_ATTR_DEVICE_BATTERY_LEVEL              = uint16(0x0015)
	C4_ATTR_RADIO_4_BARS                      = uint16(0x0016)

	//C4入网进程的几个状态
	C4_STATE_NULL       = byte(0) //初始化时
	C4_STATE_CONNECTING = byte(1) //为NUL时有report
	C4_STATE_ONLINE     = byte(2) //announce完成，offline后又收到报文，其中TC是新join的情况下announce完成会触发newnode的应用层事件
	C4_STATE_OFFLINE    = byte(3) //接收超时

	C4_MAX_OFFLINE_TIMEOUT = 300
)

type StNode struct {
	NodeID       uint16
	MAC          uint64
	LastRecvTime time.Time
	State        byte
	Newjoin      bool
	ToBeDeleted  bool

	FirmwareVersion              string
	PS                           string
	ReflashVersion               byte
	BootCount                    uint16
	DeviceType                   byte
	AnnounceWindow               uint16
	MTORRPeriod                  uint16
	EndNodeAccessPointPollPeriod uint16
	Channel                      byte
	RSSI                         int8
	LQI                          byte
}

var Nodes sync.Map

func findNodeIDbyMAC(mac uint64) (nodeID uint16) {
	nodeID = ezsp.EMBER_NULL_NODE_ID
	Nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(StNode); ok {
			if node.MAC == mac {
				nodeID = node.NodeID
				return false
			}
		}
		return true
	})
	return
}

func C4Tick() {
	select {
	case cb := <-ezsp.CallbackCh:
		ezsp.EzspCallbackDispatch(cb)
	case <-time.After(time.Millisecond * 3000):
		Nodes.Range(func(key, value interface{}) bool {
			if node, ok := value.(StNode); ok {
				node.RefreshHandle()
			}
			return true
		})
	}
}

func (node *StNode) AttribReportedHandle(z *zcl.ZclContext, cluster uint16, list []*zcl.StAttrib) error {
	var ok bool
	// todo check null
	//if z == nil || z.Context == nil {
	//	return ErrInternalError
	//}
	//dev, ok := z.Context.(models.Device)
	//if !ok {
	//	return ErrZclContextTypeMismatch
	//}

	for _, attrib := range list {
		switch attrib.AttributeIdentifier {
		case C4_ATTR_DEVICE_TYPE:
			if node.DeviceType, ok = attrib.AttributeData.(byte); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("DEVICE_TYPE: %d", node.DeviceType)
			}
		case C4_ATTR_SSCP_ANNOUNCE_WINDOW:
			if node.AnnounceWindow, ok = attrib.AttributeData.(uint16); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("ANNOUNCE_WINDOW: %d", node.AnnounceWindow)
			}
		case C4_ATTR_SSCP_MTORR_PERIOD:
			if node.MTORRPeriod, ok = attrib.AttributeData.(uint16); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("MTORR_PERIOD: %d", node.MTORRPeriod)
			}
		case C4_ATTR_FIRMWARE_VERSION:
			if node.FirmwareVersion, ok = attrib.AttributeData.(string); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("FIRMWARE_VERSION: %s", node.FirmwareVersion)
			}
		case C4_ATTR_REFLASH_VERSION:
			if node.ReflashVersion, ok = attrib.AttributeData.(byte); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("REFLASH_VERSION: %d", node.ReflashVersion)
			}
		case C4_ATTR_BOOT_COUNT:
			if node.BootCount, ok = attrib.AttributeData.(uint16); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("BOOT_COUNT: %d", node.BootCount)
			}
		case C4_ATTR_PRODUCT_STRING:
			if node.PS, ok = attrib.AttributeData.(string); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("PRODUCT_STRING: %s", node.PS)
			}
		case C4_ATTR_END_NODE_ACCESS_POINT_POLL_PERIOD:
			if node.EndNodeAccessPointPollPeriod, ok = attrib.AttributeData.(uint16); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("END_NODE_ACCESS_POINT_POLL_PERIOD: %d", node.EndNodeAccessPointPollPeriod)
			}
		case C4_ATTR_SSCP_CHANNEL:
			if node.Channel, ok = attrib.AttributeData.(byte); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("CHANNEL: %d", node.Channel)
			}
		case C4_ATTR_AVG_RSSI:
			if node.RSSI, ok = attrib.AttributeData.(int8); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("AVG_RSSI: %d", node.RSSI)
			}
		case C4_ATTR_AVG_LQI:
			if node.LQI, ok = attrib.AttributeData.(byte); !ok {
				common.Log.Errorf("Clust 0x%x attrib 0x%x parse error", cluster, attrib.AttributeIdentifier)
			} else {
				common.Log.Debugf("AVG_LQI: %d", node.LQI)
			}
		case C4_ATTR_SSCP_NUM_ZAPS:
		case C4_ATTR_ACCESS_POINT_NODE_ID:
		case C4_ATTR_ACCESS_POINT_LONG_ID:
		case C4_ATTR_ACCESS_POINT_COST:
		case C4_ATTR_ZSERVER_IP:
		case C4_ATTR_ZSERVER_HOST_NAME:
		case C4_ATTR_END_NODE_PARENT_POLL_PERIOD:
		case C4_ATTR_MESH_LONG_ID:
		case C4_ATTR_MESH_SHORT_ID:
		case C4_ATTR_AP_TABLE:
		case C4_ATTR_DEVICE_BATTERY_LEVEL:
		case C4_ATTR_RADIO_4_BARS:
		default:
		}

	}
	return nil
}

func (_ *StNode) UnsupportClusterCommandHandle(z *zcl.ZclContext, cluster uint16, direction bool, disableDefaultResponse bool, sequenceNumber byte,
	commandIdentifier byte, data interface{}) (resp []byte, err error) {
	return nil, nil
}

func (node *StNode) getState() byte {
	now := time.Now()
	if node.PS == "" && (node.DeviceType < ezsp.EMBER_ROUTER || node.DeviceType > ezsp.EMBER_MOBILE_END_DEVICE) {
		// PS未上传，且DeviceType未上传或非法
		return C4_STATE_NULL
	} else if node.PS != "" && node.DeviceType >= ezsp.EMBER_ROUTER && node.DeviceType <= ezsp.EMBER_MOBILE_END_DEVICE {
		// PS已上传，且DeviceType合法
		var timeout time.Duration
		if node.DeviceType == ezsp.EMBER_ROUTER {
			timeout = time.Duration(node.AnnounceWindow) * 2 * time.Second
		} else {
			timeout = time.Duration(node.EndNodeAccessPointPollPeriod) * 2 * time.Second
		}
		if timeout < C4_MAX_OFFLINE_TIMEOUT*time.Second {
			timeout = C4_MAX_OFFLINE_TIMEOUT * time.Second
		}
		if now.Sub(node.LastRecvTime) > timeout {
			return C4_STATE_OFFLINE
		}
		return C4_STATE_ONLINE
	} else {
		return C4_STATE_CONNECTING
	}
}

func RemoveNode(node *StNode) {
	err := ezsp.EzspRemoveDevice(node.NodeID, node.MAC, node.MAC)
	if err != nil {
		common.Log.Errorf("EzspRemoveDevice failed: %v", err)
	}

	Nodes.Delete(node.NodeID)
}

// RefreshHandle 收到报文发生变化，或定时刷新时调用
func (node *StNode) RefreshHandle() {
	if node.ToBeDeleted {
		// todo 通知delete
		Nodes.Delete(node.NodeID)
		common.Log.Infof("node 0x%016x deleted", node.MAC)
		return
	}

	newState := node.getState()

	if newState != node.State {
		if newState == C4_STATE_ONLINE {
			if node.State < C4_STATE_ONLINE && node.Newjoin { //初次入网，且NULL、CONNECTING变成ONLINE，要检查passport，不允许的踢出
				if passportMatch(node.PS, fmt.Sprintf("%016x", node.MAC)) >= 0 {
					common.Log.Infof("%s node 0x%016x join C4 network", node.PS, node.MAC)
				} else {
					common.Log.Errorf("reject node 0x%016x, remove it", node.MAC)
					RemoveNode(node)
					return
				}
			} else {
				common.Log.Infof("node 0x%016x reonline", node.MAC)
			}
			// todo 通知online
		} else if newState == C4_STATE_OFFLINE {
			common.Log.Infof("node 0x%016x offline", node.MAC)
			// todo 通知offline
		}
		node.State = newState
	}
	Nodes.Store(node.NodeID, *node) // map中存储
}

type StPassport struct {
	PS  string
	MAC string
}

var allPassPorts []*StPassport

//返回Match的字符个数，不含x
func checkPassportMAC(mac string, ppMac string) (match int) {
	if len([]rune(mac)) != 16 {
		common.Log.Errorf("MAC len not 16: %s", mac)
		return -1
	}
	for i, c := range mac {
		ppc := []rune(ppMac)[i]
		if !(((c >= '0') && (c <= '9')) || ((c >= 'a') && (c <= 'f'))) {
			return -1
		}
		if ppc == 'x' {
			continue
		}
		if ppc != c {
			return -1
		}
		match++
	}
	return
}

func passportMatch(ps string, mac string) (maxHit int) {
	maxHit = -1
	if allPassPorts != nil {
		for _, p := range allPassPorts {
			if ps == p.PS {
				match := checkPassportMAC(mac, p.MAC)
				if match > maxHit {
					maxHit = match
					if maxHit >= 16 {
						break
					}
				}
			}
		}
	}
	return
}

//16位hex字符或前置若干个x都是合法的
func isPassportMACValid(mac string) bool {
	leadingx := true
	if len([]rune(mac)) != 16 {
		common.Log.Errorf("MAC len not 16: %s", mac)
		return false
	}
	for _, c := range mac {
		//hex字符一定正确
		if ((c >= '0') && (c <= '9')) || ((c >= 'a') && (c <= 'f')) {
			leadingx = false //出现了hex字符
			continue
		}
		//
		if (c != 'x') || (leadingx == false) {
			return false
		}
	}
	return true
}

func C4SetPermission(duration byte, passports []*StPassport) (err error) {

	if passports == nil {
		err = fmt.Errorf("C4SetPassports passports=NULL")
		return
	}
	common.Log.Debugf("permision to ...")
	for i, p := range passports {
		if p != nil {
			common.Log.Debugf("%d: MAC=%s PS=%s", i, p.MAC, p.PS)
			if p.PS == "" {
				err = fmt.Errorf("C4SetPassports passport %d PS=NULL", i)
				return
			}
			if isPassportMACValid(p.MAC) == false {
				err = fmt.Errorf("C4SetPassports passport %d mac -%s- invalid", i, p.MAC)
				return
			}
		}
	}
	allPassPorts = passports

	err = ezsp.EzspPermitJoining(duration)
	if err != nil {
		err = fmt.Errorf("EzspPermitJoining failed: %v", err)
	}

	return
}

func C4Init() {
	ezsp.Networker.NcpTrustCenterJoinHandler = TrustCenterJoinHandler
	ezsp.Networker.NcpMessageSentHandler = MessageSentHandler
	ezsp.Networker.NcpIncomingMessageHandler = IncomingMessageHandler
}

func TrustCenterJoinHandler(newNodeId uint16,
	newNodeEui64 uint64,
	deviceUpdateStatus byte,
	joinDecision byte,
	parentOfNewNode uint16) {

	now := time.Now()
	var node StNode
	value, ok := Nodes.Load(newNodeId) // 从map中加载
	if ok {
		if node, ok = value.(StNode); !ok {
			common.Log.Errorf("Nodes map unsupported type")
			return
		}
		node.Newjoin = false //已经在node表里的，不认为Newjoin
	} else {
		node = StNode{NodeID: newNodeId, MAC: newNodeEui64, LastRecvTime: now}
		node.Newjoin = true
	}

	if (deviceUpdateStatus == ezsp.EMBER_STANDARD_SECURITY_UNSECURED_JOIN) ||
		(deviceUpdateStatus == ezsp.EMBER_HIGH_SECURITY_UNSECURED_JOIN) {
	} else if (deviceUpdateStatus == ezsp.EMBER_STANDARD_SECURITY_SECURED_REJOIN) ||
		(deviceUpdateStatus == ezsp.EMBER_STANDARD_SECURITY_UNSECURED_REJOIN) ||
		(deviceUpdateStatus == ezsp.EMBER_HIGH_SECURITY_SECURED_REJOIN) ||
		(deviceUpdateStatus == ezsp.EMBER_HIGH_SECURITY_UNSECURED_REJOIN) {
		node.Newjoin = false
	} else if deviceUpdateStatus == ezsp.EMBER_DEVICE_LEFT {
		node.ToBeDeleted = true
	} else {
		common.Log.Errorf("unknown deviceUpdateStatus = 0x%x newNodeId = 0x%04x\r\n", deviceUpdateStatus, newNodeId)
		return
	}
	node.RefreshHandle()
}

func MessageSentHandler(outgoingMessageType byte,
	indexOrDestination uint16,
	apsFrame *ezsp.EmberApsFrame,
	messageTag byte,
	emberStatus byte,
	message []byte) {

}
func IncomingMessageHandler(incomingMessageType byte,
	apsFrame *ezsp.EmberApsFrame,
	lastHopLqi byte,
	lastHopRssi int8,
	sender uint16,
	bindingIndex byte,
	addressIndex byte,
	message []byte) {

	now := time.Now()
	eui64, err := ezsp.EzspLookupEui64ByNodeId(sender)
	if err != nil {
		common.Log.Errorf("Incoming message lookup eui64 failed: %v", err)
		return
	}

	var node StNode
	value, ok := Nodes.Load(sender) // 从map中加载
	if ok {
		if node, ok = value.(StNode); !ok {
			common.Log.Errorf("Nodes map unsupported type")
			return
		}
	} else {
		node = StNode{NodeID: sender, MAC: eui64, LastRecvTime: now}
	}

	if apsFrame.ProfileId == C4_PROFILE {
		if apsFrame.ClusterId == C4_CLUSTER {
			zclContext := &zcl.ZclContext{LocalEdp: apsFrame.DestinationEndpoint, RemoteEdp: apsFrame.SourceEndpoint,
				Context:      nil,
				GlobalHandle: &node}

			resp, err := zclContext.Parse(apsFrame.ProfileId, apsFrame.ClusterId, message)
			if err != nil {
				common.Log.Errorf("Incoming C25D message parse failed: %v", err)
				return
			}

			if incomingMessageType == ezsp.EMBER_INCOMING_UNICAST && resp != nil && len(resp) > 0 {
				var respFrame ezsp.EmberApsFrame
				respFrame.ProfileId = apsFrame.ProfileId
				respFrame.ClusterId = apsFrame.ClusterId
				respFrame.SourceEndpoint = apsFrame.DestinationEndpoint
				respFrame.DestinationEndpoint = apsFrame.SourceEndpoint
				respFrame.Options = getSendOptions(sender, respFrame.ProfileId, respFrame.ClusterId, byte(len(resp)))

				ezsp.EzspSendUnicast(ezsp.EMBER_OUTGOING_DIRECT, sender, &respFrame, 0, resp)
			} else if incomingMessageType == ezsp.EMBER_INCOMING_BROADCAST {
				ezsp.NcpSendMTORR()
			}

			node.RefreshHandle()
		}

	} else if apsFrame.ProfileId == ZDO_PROFILE {

	} else {
		//todo 通知
	}
}

var unicastTagSequence = byte(0)

func nextSequence() byte {
	unicastTagSequence++
	if unicastTagSequence == 0 || unicastTagSequence == 0xff {
		unicastTagSequence = 1
	}
	return unicastTagSequence
}

func C4SendUnicast(mac uint64,
	profileId uint16,
	clusterId uint16,
	localEndpoint byte,
	remoteEndpoint byte,
	message []byte,
	needConfirm bool) {

	nodeID := findNodeIDbyMAC(mac)
	var apsFrame ezsp.EmberApsFrame
	apsFrame.ProfileId = profileId
	apsFrame.ClusterId = clusterId
	apsFrame.SourceEndpoint = localEndpoint
	apsFrame.DestinationEndpoint = remoteEndpoint
	apsFrame.Options = getSendOptions(nodeID, profileId, clusterId, byte(len(message)))
	tag := byte(0)
	if needConfirm {
		tag = nextSequence()
	}

	ezsp.EzspSendUnicast(ezsp.EMBER_OUTGOING_DIRECT, nodeID, &apsFrame, tag, message)
}

func getSendOptions(destination uint16, profileId uint16, clusterId uint16, messageLength byte) (options uint16) {
	if profileId == 0xc25d && clusterId == 0x0001 {
		if destination >= ezsp.EMBER_BROADCAST_ADDRESS {
			options = /*ezsp.EMBER_APS_OPTION_RETRY |*/ ezsp.EMBER_APS_OPTION_SOURCE_EUI64
		} else {
			options = ezsp.EMBER_APS_OPTION_RETRY | ezsp.EMBER_APS_OPTION_ENABLE_ROUTE_DISCOVERY | ezsp.EMBER_APS_OPTION_ENABLE_ADDRESS_DISCOVERY | ezsp.EMBER_APS_OPTION_SOURCE_EUI64 | ezsp.EMBER_APS_OPTION_DESTINATION_EUI64
		}
	} else {

		if messageLength <= 66 { /*66不溢*/
			options = /*ezsp.EMBER_APS_OPTION_RETRY | */ ezsp.EMBER_APS_OPTION_SOURCE_EUI64 | ezsp.EMBER_APS_OPTION_DESTINATION_EUI64
		} else if messageLength <= 74 { /*67~74*/
			options = /*ezsp.EMBER_APS_OPTION_RETRY | */ ezsp.EMBER_APS_OPTION_SOURCE_EUI64
		} else {
			options = ezsp.EMBER_APS_OPTION_NONE
		}
	}
	return
}
