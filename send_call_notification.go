package main

import nex "github.com/PretendoNetwork/nex-go"

func sendCallNotification(caller uint32, target uint32, callID uint32) {
	event := &NotificationEvent{}

	event.sourcePID = caller                 // Sender PID
	event.typeParameter = 101000             // Notification type
	event.parameter1 = caller                // Gathering ID
	event.parameter2 = target                // Recipient PID
	event.stringParameter = "Invite Request" // Unknown

	eventObject := nex.NewStreamOut(nexServer)
	eventObject.WriteStructure(event)

	rmcRequest, _ := nex.NewRMCRequest([]byte{})
	rmcRequest.SetProtocolID(14)
	rmcRequest.SetCallID(0xffff + callID)
	rmcRequest.SetMethodID(1)
	rmcRequest.SetParameters(eventObject.Bytes())

	rmcRequestBytes := rmcRequest.Bytes()

	clientAddr := getPlayerSessionAddress(target)

	requestPacket, _ := nex.NewPacketV1(nexServer.GetClient(clientAddr), nil)

	requestPacket.SetVersion(1)
	requestPacket.SetSource(0xA1)
	requestPacket.SetDestination(0xAF)
	requestPacket.SetType(nex.DataPacket)
	requestPacket.SetPayload(rmcRequestBytes)

	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)

	nexServer.Send(requestPacket)
}
