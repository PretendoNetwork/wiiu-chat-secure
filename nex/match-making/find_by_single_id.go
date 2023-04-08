package nex_match_making

import (
	nex "github.com/PretendoNetwork/nex-go"
	match_making "github.com/PretendoNetwork/nex-protocols-go/match-making"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func FindBySingleID(err error, client *nex.Client, callID uint32, id uint32) {
	caller := id // Gathering ID and caller are the same here

	result := true

	gathering := match_making.NewGathering()
	gathering.ID = id
	gathering.OwnerPID = caller
	gathering.HostPID = caller
	gathering.MinimumParticipants = 2
	gathering.MaximumParticipants = 2
	gathering.Description = "Doors Invite Request"

	dataHolder := nex.NewDataHolder()
	dataHolder.SetTypeName("Gathering")
	dataHolder.SetObjectData(gathering)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	rmcResponseStream.WriteBool(result)
	rmcResponseStream.WriteDataHolder(dataHolder)

	rmcResponse := nex.NewRMCResponse(match_making.ProtocolID, callID)
	rmcResponse.SetSuccess(match_making.MethodFindBySingleID, rmcResponseStream.Bytes())

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
