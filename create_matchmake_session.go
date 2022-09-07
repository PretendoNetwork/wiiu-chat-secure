package main

import (
	//"math/rand"

	nex "github.com/PretendoNetwork/nex-go"

	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func createMatchmakeSession(err error, client *nex.Client, callID uint32, matchmakeSession *nexproto.MatchmakeSession, message string, participationCount uint16) {
	var gid uint32 = client.PID() // TODO: Random this
	sessionKey := make([]byte, 32)
	//rand.Read(sessionKey)

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt32LE(gid)
	rmcResponseStream.WriteBuffer(sessionKey)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodCreateMatchmakeSession, rmcResponseBody)

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
