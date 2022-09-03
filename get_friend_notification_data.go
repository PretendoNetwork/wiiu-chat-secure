package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"

	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getFriendNotificationData(err error, client *nex.Client, callID uint32, uiType int32) {

	fmt.Printf("uiType: %d\r\n", uiType)

	rmcResponseStream := nex.NewStreamOut(nexServer)

	/*
		// List<NotificationEvent>

		// This enableds auto-match making for calls
		var caller uint32 = 1743126339
		var target uint32 = 1424784406

		event := nexproto.NewNotificationEvent()

		event.PIDSource = caller          // Sender PID
		event.Type = 101000               // Notification type
		event.Param1 = caller             // Gathering ID
		event.Param2 = target             // Recipient PID
		event.StrParam = "Invite Request" // Unknown

		rmcResponseStream.WriteUInt32LE(1)
		rmcResponseStream.WriteStructure(event)
	*/

	rmcResponseStream.WriteUInt32LE(0) // No data for now

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodGetFriendNotificationData, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)

	//// HANDLE INCOMING CALL ////

	caller, target, ringing := getCallInfoByTarget(client.PID())
	if (caller != 0) && (target == client.PID()) && ringing {
		sendCallNotification(caller, target, callID)
	}
}
