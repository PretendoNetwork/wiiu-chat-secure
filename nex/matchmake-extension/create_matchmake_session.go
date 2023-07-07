package nex_matchmake_extension

import (
	//"math/rand"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
)

func CreateMatchmakeSession(err error, client *nex.Client, callID uint32, anyGathering *nex.DataHolder, message string, participationCount uint16) {
	var gid uint32 = client.PID() // TODO: Random this
	sessionKey := make([]byte, 32)
	//rand.Read(sessionKey)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteUInt32LE(gid)
	rmcResponseStream.WriteBuffer(sessionKey)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodCreateMatchmakeSession, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
