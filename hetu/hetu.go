package hetu

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/conthing/ezsp/ezsp"

	"encoding/binary"

	"github.com/conthing/utils/common"
	"github.com/conthing/utils/crc16"
)

const (
	ZDO_PROFILE = uint16(0x0000)

	//C4入网进程的几个状态
	C4_STATE_NULL       = byte(0) //初始化时
	C4_STATE_CONNECTING = byte(1) //为NUL时有report
	C4_STATE_ONLINE     = byte(2) //announce完成，offline后又收到报文，其中TC是新join的情况下announce完成会触发newnode的应用层事件
	C4_STATE_OFFLINE    = byte(3) //接收超时

	C4_MAX_OFFLINE_TIMEOUT = 90 //未收到报文达到90秒，标记为offline
)

var ZeroTime = time.Unix(0, 0)           // 1970-1-1 00:00:00 既没announce也没有收到有效报文
var AnnouceFlagTime = time.Unix(3661, 0) // 1970-1-1 01:01:01 表示设备已经announce但还没有收到有效报文的lastrecvtime

var ErrMeshNotExist = errors.New("Mesh not exist")
var ErrMeshAlreadyExist = errors.New("Mesh already exist")
var ErrMeshNotEmpty = errors.New("Not empty mesh")

type StNode struct {
	Eui64        uint64
	NodeID       uint16
	State        byte
	Addr         byte
	LastRecvTime time.Time
}

//eui64 要转成 16进制 mac
type StC4Callbacks struct {
	C4MessageSentHandler     func(eui64 uint64, profileId uint16, clusterId uint16, localEndpoint byte, remoteEndpoint byte, message []byte, success bool)
	C4IncomingMessageHandler func(eui64 uint64, message []byte, recvTime time.Time)
	C4NodeStatusHandler      func(eui64 uint64, nodeID uint16, status byte, addr byte)
}

var C4Callbacks StC4Callbacks

var Nodes sync.Map

// LoadNodesMap 加载 Map
func LoadNodesMap(m map[uint64]StNode) {
	for _, node := range m {
		StoreNode(&node)
		common.Log.Info("1 LoadNodesMap: ", node.NodeID)
	}
}

// StoreNode 保存节点到Nodes中，如果有重复的eui64，更新
func StoreNode(node *StNode) {
	nodeID := findNodeIDbyEui64(node.Eui64)
	if nodeID == ezsp.EMBER_NULL_NODE_ID {
		Nodes.Store(node.NodeID, *node) // map中存储
	} else {
		Nodes.Delete(nodeID)            // map中原来的删掉
		Nodes.Store(node.NodeID, *node) // map中存储
	}
}

//在Nodes中找到匹配的eui64
func findNodeIDbyEui64(eui64 uint64) (nodeID uint16) {
	nodeID = ezsp.EMBER_NULL_NODE_ID
	Nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(StNode); ok {
			if node.Eui64 == eui64 {
				nodeID = node.NodeID
				return false
			}
		}
		return true
	})
	return
}

var lastTimeStamp = int64(0)

func HetuTick() {
	select {
	case cbs := <-ezsp.CallbackCh:
		for _, cb := range cbs {
			ezsp.EzspCallbackDispatch(cb)
		}
	case <-time.After(time.Millisecond * 500):

	}
	now := time.Now().Unix()
	if now-lastTimeStamp >= 10 {
		lastTimeStamp = now
		Nodes.Range(func(key, value interface{}) bool {
			if node, ok := value.(StNode); ok {
				node.RefreshHandle(false)
			}
			return true
		})
		if ezsp.MeshStatusUp {
			err := HetuBroadcast()
			if err != nil {
				common.Log.Errorf("hetu tick: %v", err)
			}
		}
	}
}

func (node *StNode) getState() byte {
	now := time.Now()

	if node.LastRecvTime == AnnouceFlagTime { //announce但没有报文上来
		return C4_STATE_CONNECTING
	}

	timeout := C4_MAX_OFFLINE_TIMEOUT * time.Second
	if now.Sub(node.LastRecvTime) > timeout {
		return C4_STATE_OFFLINE
	}
	return C4_STATE_ONLINE
}

func removeDeviceAndNode(node *StNode) {
	err := ezsp.EzspRemoveDevice(node.NodeID, node.Eui64, node.Eui64)
	if err != nil {
		common.Log.Errorf("EzspRemoveDevice failed: %v", err)
	}

	Nodes.Delete(node.NodeID)
}

// RefreshHandle 收到报文发生变化，或定时刷新时调用
func (node *StNode) RefreshHandle(forceReport bool) {

	newState := node.getState()
	common.Log.Debugf("newState:%d", newState)

	if newState != node.State || forceReport {
		if newState == C4_STATE_ONLINE {
			if node.State < C4_STATE_ONLINE { //初次入网，且NULL、CONNECTING变成ONLINE，要检查passport，不允许的踢出
				common.Log.Infof("node 0x%016x online", node.Eui64)
			} else {
				common.Log.Infof("node 0x%016x reonline", node.Eui64)
			}
			common.Log.Debugf("HetuNodeStatusHandler online")
			if C4Callbacks.C4NodeStatusHandler != nil {
				C4Callbacks.C4NodeStatusHandler(node.Eui64, node.NodeID, C4_STATE_ONLINE, node.Addr)
			}
		} else if newState == C4_STATE_OFFLINE {
			common.Log.Infof("node 0x%016x offline", node.Eui64)
			common.Log.Debugf("HetuNodeStatusHandler offline")
			if C4Callbacks.C4NodeStatusHandler != nil {
				C4Callbacks.C4NodeStatusHandler(node.Eui64, node.NodeID, C4_STATE_OFFLINE, node.Addr)
			}
		}
		node.State = newState
	}
	Nodes.Store(node.NodeID, *node) // map中存储
}

func SetPermission(duration byte) (err error) {
	common.Log.Debugf("Permit join for %d seconds", duration)

	err = ezsp.EzspPermitJoining(duration)
	if err != nil {
		err = fmt.Errorf("EzspPermitJoining failed: %v", err)
	}
	return
}

func Init() {
	ezsp.NcpCallbacks.NcpMessageSentHandler = MessageSentHandler
	ezsp.NcpCallbacks.NcpIncomingMessageHandler = IncomingMessageHandler
}

var hndl_cnt byte

func HetuBroadcast() error {
	apsFrame := ezsp.EmberApsFrame{ProfileId: 0xabcd, ClusterId: 0xabef, SourceEndpoint: 2, DestinationEndpoint: 2}
	hndl_cnt++
	if hndl_cnt >= 0xfe {
		hndl_cnt = 0
	}
	message := []byte{0x78, 0x87, hndl_cnt}
	_, err := ezsp.EzspSendBroadcast(ezsp.EMBER_SLEEPY_BROADCAST_ADDRESS, &apsFrame, 30, 0, message)
	return err
}

func MessageSentHandler(outgoingMessageType byte,
	indexOrDestination uint16,
	apsFrame *ezsp.EmberApsFrame,
	messageTag byte,
	emberStatus byte,
	message []byte) {
	if apsFrame.ProfileId == ZDO_PROFILE {
		return
	}
	if messageTag == 0 { //应用层不关心时tag==0
		return
	}
	var node StNode
	value, ok := Nodes.Load(indexOrDestination) // 从map中加载
	if ok {
		if node, ok = value.(StNode); !ok {
			common.Log.Errorf("Nodes map unsupported type")
			return
		}
	} else {
		common.Log.Errorf("0x%04x not found in Nodes map", indexOrDestination)
		return
	}
	if C4Callbacks.C4MessageSentHandler != nil {
		C4Callbacks.C4MessageSentHandler(node.Eui64, apsFrame.ProfileId, apsFrame.ClusterId, apsFrame.SourceEndpoint, apsFrame.DestinationEndpoint, message, emberStatus == ezsp.EMBER_SUCCESS)
	}
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
	if apsFrame.ProfileId == ZDO_PROFILE {
		if apsFrame.ClusterId == 0x0013 { //device announce
			nodeID := binary.LittleEndian.Uint16(message[1:])
			eui64 := binary.LittleEndian.Uint64(message[3:])
			node := StNode{NodeID: nodeID, Eui64: eui64, LastRecvTime: AnnouceFlagTime}
			StoreNode(&node)
			node.RefreshHandle(false)
			common.Log.Debugf("2 zdo announce: 0x%04x,%016x", nodeID, eui64)
		}
	} else {
		if incomingMessageType == ezsp.EMBER_INCOMING_UNICAST {
			if apsFrame.ProfileId == 0xabcd && apsFrame.ClusterId == 0xabde && (len(message) == 38 || len(message) == 70) && crc16.CRC16MODBUS(message) == 0 {
				// 收到一条合法的hetu报文
				var forceReport bool
				var node StNode
				value, ok := Nodes.Load(sender) // 从map中加载
				if ok {
					if node, ok = value.(StNode); !ok {
						common.Log.Errorf("Nodes map unsupported type")
						return
					}
					common.Log.Debugf("Nodes map get %016x", node.Eui64)
					if node.Addr != message[5] {
						if node.LastRecvTime != AnnouceFlagTime && node.LastRecvTime != ZeroTime {
							common.Log.Warnf("Node %016x digital addr changed from %d to %d", node.Eui64, node.Addr, message[5])
						}
						node.Addr = message[5]
						forceReport = true
					}
					node.LastRecvTime = now
				} else {
					eui64, err := ezsp.EzspLookupEui64ByNodeId(sender)
					if err != nil {
						common.Log.Errorf("Incoming message lookup eui64 failed: %v", err)
						return
					}

					common.Log.Debugf("EzspLookupEui64ByNodeId get %016x", eui64)
					node = StNode{NodeID: sender, Eui64: eui64, Addr: message[5], LastRecvTime: now}
				}
				StoreNode(&node)
				node.RefreshHandle(forceReport)
				common.Log.Debugf("3 HetuIncomingMessageHandler: %d", node.NodeID)
				if C4Callbacks.C4IncomingMessageHandler != nil {
					if node.Eui64 != 0 {
						C4Callbacks.C4IncomingMessageHandler(node.Eui64, message, now)
					} else {
						common.Log.Errorf("recv msg from NodeID 0x%04x without EUI64", node.NodeID)
					}
				}
			} else {
				common.Log.Errorf("Incoming invalid message Profile=0x%04x Cluster=0x%04x CRC=0x%04x", apsFrame.ProfileId, apsFrame.ClusterId, crc16.CRC16MODBUS(message))
			}
		}
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

func SendUnicast(eui64 uint64, profileId uint16, clusterId uint16,
	localEndpoint byte, remoteEndpoint byte,
	message []byte, needConfirm bool) (err error) {
	common.Log.Debugf("SendUnicast %016x ...", eui64)

	nodeID := findNodeIDbyEui64(eui64)
	if nodeID == ezsp.EMBER_NULL_NODE_ID {
		return fmt.Errorf("unknow EUI64 %016x", eui64)
	}
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

	// todo 设置路由表
	_ = ezsp.NcpSetSourceRoute(nodeID)
	_, err = ezsp.EzspSendUnicast(ezsp.EMBER_OUTGOING_DIRECT, nodeID, &apsFrame, tag, message)
	return
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

func RemoveDevice(eui64 uint64) (err error) {
	common.Log.Debugf("RemoveDevice %016x", eui64)
	nodeID := findNodeIDbyEui64(eui64)
	if nodeID == ezsp.EMBER_NULL_NODE_ID {
		return fmt.Errorf("unknow EUI64 %016x", eui64)
	}
	err = ezsp.EzspRemoveDevice(nodeID, eui64, eui64)
	if err != nil {
		common.Log.Errorf("EzspRemoveDevice failed: %v", err)
	}
	return
}

func RemoveNetwork() (err error) {
	common.Log.Debugf("RemoveNetwork()")
	if !ezsp.MeshStatusUp {
		return ErrMeshNotExist
	}

	Nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(StNode); ok {
			removeDeviceAndNode(&node)
		}
		return true
	})

	err = ezsp.EzspLeaveNetwork()
	return
}

func SetRadioChannel(channel byte) (err error) {
	common.Log.Debugf("SetRadioChannel(%d)", channel)
	return ezsp.EzspSetRadioChannel(channel)
}

func FormNetwork(radioChannel byte) (err error) {
	common.Log.Debugf("FormNetwork(%d)", radioChannel)
	if ezsp.MeshStatusUp {
		return ErrMeshAlreadyExist
	} else {
		return ezsp.NcpFormNetwork(radioChannel, false)
	}
}
