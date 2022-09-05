package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func findBySingleID(err error, client *nex.Client, callID uint32, id uint32) {
	caller := id // Gathering ID and caller are the same here

	result := true

	gathering := nexproto.NewGathering()
	gathering.ID = id
	gathering.OwnerPID = caller
	gathering.HostPID = caller
	gathering.MinimumParticipants = 2
	gathering.MaximumParticipants = 2
	gathering.Description = "Doors Invite Request"

	dataHolder := nex.NewDataHolder()
	dataHolder.SetTypeName("Gathering")
	dataHolder.SetObjectData(gathering)

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteBool(result)
	rmcResponseStream.WriteDataHolder(dataHolder)

	rmcResponse := nex.NewRMCResponse(nexproto.MatchMakingMethodFindBySingleID, callID)
	rmcResponse.SetSuccess(nexproto.MatchMakingProtocolID, rmcResponseStream.Bytes())

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
