package main

import (
	nex "github.com/PretendoNetwork/nex-go"

	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getFriendNotificationData(err error, client *nex.Client, callID uint32, uiType int32) {
	notifications := make([]*nexproto.NotificationEvent, 0)

	caller, target, ringing := getCallInfoByTarget(client.PID())

	// TODO: Multiple calls. Wii U Chat can handle it, but we don't support it yet
	if (caller != 0) && (target == client.PID()) && ringing {
		// Being called
		notification := nexproto.NewNotificationEvent()

		notification.PIDSource = caller
		notification.Type = 101000
		notification.Param1 = caller
		notification.Param2 = target
		notification.StrParam = "Invite Request"

		notifications = append(notifications, notification)
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteListStructure(notifications)

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
}
