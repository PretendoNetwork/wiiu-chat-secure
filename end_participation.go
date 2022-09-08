package main

import (
	"strconv"

	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func endParticipation(err error, client *nex.Client, callID uint32, idGathering uint32, strMessage string) {
	logger.Info("ENDING PARTICIPATION: GID: " + strconv.FormatUint(uint64(idGathering), 10) + ", MESSAGE: " + strMessage)

	rmcResponse := nex.NewRMCResponse(nexproto.MatchMakingExtMethodEndParticipation, callID)
	rmcResponse.SetSuccess(nexproto.MatchMakingExtProtocolID, []byte{0x01})

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
