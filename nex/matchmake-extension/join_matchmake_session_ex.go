package nex_matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func JoinMatchmakeSessionEx(err error, client *nex.Client, callID uint32, gid uint32, strMessage string, dontCareMyBlockList bool, participationCount uint16) {
	globals.Logger.Infof("gid: %d, strMessage: %s, dontCareMyBlockList: %t, participationCount: %d\r\n", gid, strMessage, dontCareMyBlockList, participationCount)

	database.EndCallRinging(gid)

	output := nex.NewStreamOut(globals.NEXServer)
	output.WriteBuffer(make([]byte, 32))

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodJoinMatchmakeSessionEx, output.Bytes())

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
