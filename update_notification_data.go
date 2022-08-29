package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"

	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func updateNotificationData(err error, client *nex.Client, callID uint32, uiType uint32, uiParam1 uint32, uiParam2 uint32, strParam string) {
	// TODO: Implement this
	fmt.Printf("uiType: %d, uiParam1: %d, uiParam2: %d, strParam: %s\r\n", uiType, uiParam1, uiParam2, strParam)

	// kick player if invite cancellation to prevent app hanging indefinitely
	if uiType == 102 {
		endCall(uiParam1)
		nexServer.Kick(client)
		return
	}

	if uiType == 101 {
		newCall(client.PID(), uiParam2)
		if doesSessionExist(uiParam2) {
			sendCallNotification(client.PID(), uiParam2, callID)
		}
	}

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodUpdateNotificationData, nil)

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
