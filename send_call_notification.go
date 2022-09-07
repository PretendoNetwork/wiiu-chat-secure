package main

import (
	//"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func sendCallNotification(caller uint32, target uint32, callID uint32) {
	event := nexproto.NewNotificationEvent()

	event.PIDSource = caller          // Sender PID
	event.Type = 101000               // Notification type
	event.Param1 = caller             // Gathering ID
	event.Param2 = target             // Recipient PID
	event.StrParam = "Invite Request" // Unknown

	eventObject := nex.NewStreamOut(nexServer)
	eventObject.WriteStructure(event)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(14)
	rmcRequest.SetCallID(0xffff + callID)
	rmcRequest.SetMethodID(1)
	rmcRequest.SetParameters(eventObject.Bytes())

	rmcRequestBytes := rmcRequest.Bytes()

	targetClient := nexServer.FindClientFromPID(uint32(target))

	requestPacket, _ := nex.NewPacketV1(targetClient, nil)

	requestPacket.SetVersion(1)
	requestPacket.SetSource(0xA1)
	requestPacket.SetDestination(0xAF)
	requestPacket.SetType(nex.DataPacket)
	requestPacket.SetPayload(rmcRequestBytes)

	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)

	nexServer.Send(requestPacket)
}
