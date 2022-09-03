package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func findBySingleID(err error, client *nex.Client, callID uint32, id uint32) {
	caller, _, _ := getCallInfoByCaller(id)

	gathering := nexproto.NewGathering()
	gathering.ID = id
	gathering.OwnerPID = caller
	gathering.HostPID = caller
	gathering.MinimumParticipants = 2
	gathering.MaximumParticipants = 2
	gathering.Description = "Doors Invite Request"

	outStream := nex.NewStreamOut(nexServer)
	outStream.WriteBool(true)
	outStream.WriteString("Gathering")
	b := gathering.Bytes(nex.NewStreamOut(nexServer))
	outStream.WriteUInt32LE(uint32(4) + uint32(len(b)))
	outStream.WriteBuffer(b)

	rmcResponse := nex.NewRMCResponse(nexproto.MatchMakingMethodFindBySingleID, callID)
	rmcResponse.SetSuccess(nexproto.MatchMakingProtocolID, outStream.Bytes())

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
