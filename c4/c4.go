package c4
import (
	ezsp "github.com/conthing/ezsp/ezsp"
)

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

}

