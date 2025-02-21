package nex_match_making

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func GetSessionUrls(err error, packet nex.PacketInterface, callID uint32, gid types.UInt32) (*nex.RMCMessage, *nex.Error) {
	caller := types.NewPID(uint64(gid)) // Gathering ID and caller are the same here

	stationUrlStrings := globals.SecureEndpoint.FindConnectionByPID(uint64(caller)).StationURLs

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	stationUrlStrings.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = match_making.ProtocolID
	rmcResponse.CallID = callID
	rmcResponse.MethodID = match_making.MethodGetSessionURLs

	return rmcResponse, nil
}
