package nex_match_making_ext

import (
	"strconv"

	"github.com/PretendoNetwork/nex-go"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/match-making-ext"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func EndParticipation(err error, client *nex.Client, callID uint32, idGathering uint32, strMessage string) {
	globals.Logger.Info("ENDING PARTICIPATION: GID: " + strconv.FormatUint(uint64(idGathering), 10) + ", MESSAGE: " + strMessage)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	// TODO - Check if this is always true or not, and if not when to set to false
	rmcResponseStream.WriteBool(true)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(match_making_ext.ProtocolID, callID)
	rmcResponse.SetSuccess(match_making_ext.MethodEndParticipation, rmcResponseBody)

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
