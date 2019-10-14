package c4

import (
	"sync"
	"time"

	ezsp "github.com/conthing/ezsp/ezsp"
	"github.com/conthing/ezsp/zcl"
	"github.com/conthing/utils/common"
)

type StNode struct {
	NodeID       uint16
	MAC          uint64
	LastRecvTime time.Time
}

var Nodes sync.Map

func C4Init() {
	ezsp.Networker.NcpMessageSentHandler = C4MessageSentHandler
	ezsp.Networker.NcpIncomingSenderEui64Handler = C4IncomingSenderEui64Handler
	ezsp.Networker.NcpIncomingMessageHandler = C4IncomingMessageHandler

}

func C4MessageSentHandler(outgoingMessageType byte,
	indexOrDestination uint16,
	apsFrame *ezsp.EmberApsFrame,
	messageTag byte,
	emberStatus byte,
	message []byte) {

}
func C4IncomingSenderEui64Handler(senderEui64 uint64) {

}
func C4IncomingMessageHandler(incomingMessageType byte,
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
	Nodes.Store(sender, StNode{NodeID: sender, MAC: eui64, LastRecvTime: now})

	if apsFrame.ProfileId == 0xc25d {
		zclContext := &zcl.ZclContext{LocalEdp: apsFrame.DestinationEndpoint, RemoteEdp: apsFrame.SourceEndpoint,
			Context:      nil,
			GlobalHandle: nil}

		resp, err := zclContext.Parse(apsFrame.ProfileId, apsFrame.ClusterId, message)
		if err != nil {
			common.Log.Errorf("Incoming C25D message parse failed: %v", err)
			return
		}
		if resp != nil && len(resp) > 0 {
			//data := znet_models.StUnicast{}
			//data.NeedConfirm = false
			//data.Msg.MAC = recv.Msg.MAC
			//data.Msg.Profile = recv.Msg.Profile
			//data.Msg.Cluster = recv.Msg.Cluster
			//data.Msg.LocalEdp = recv.Msg.LocalEdp
			//data.Msg.RemoteEdp = recv.Msg.RemoteEdp
			//data.Msg.Data = resp
			//
			//SendUnicast(data)
		}
	}
}
