package nex_match_making

import (
	"github.com/PretendoNetwork/nex-go/v2/types"

	nex "github.com/PretendoNetwork/nex-go/v2"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func FindBySingleID(err error, packet nex.PacketInterface, callID uint32, id types.UInt32) (*nex.RMCMessage, *nex.Error) {
	caller := types.NewPID(uint64(id)) // Gathering ID and caller are the same here

	result := types.NewBool(true)

	gathering := match_making_types.NewGathering()
	gathering.ID = id
	gathering.OwnerPID = caller
	gathering.HostPID = caller
	gathering.MinimumParticipants = 2
	gathering.MaximumParticipants = 2
	gathering.Description = "Doors Invite Request"

	// Gathering does not support being within an AnyDataHolder for some reason,
	// so we have to create the AnyDataHolder info manually
	dataHolder := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	gathering.WriteTo(dataHolder)
	dataHolderBytes := dataHolder.Bytes()

	name := types.NewString("Gathering")
	holderLen := types.NewUInt32(uint32(len(dataHolderBytes)) + 4)
	holderBuf := types.NewBuffer(dataHolderBytes)

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	result.WriteTo(rmcResponseStream)
	name.WriteTo(rmcResponseStream)
	holderLen.WriteTo(rmcResponseStream)
	holderBuf.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = match_making.ProtocolID
	rmcResponse.CallID = callID
	rmcResponse.MethodID = match_making.MethodFindBySingleID

	return rmcResponse, nil
}
