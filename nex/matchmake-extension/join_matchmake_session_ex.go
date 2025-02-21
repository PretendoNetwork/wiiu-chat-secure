package nex_matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	"github.com/PretendoNetwork/wiiu-chat/database"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func JoinMatchmakeSessionEx(err error, packet nex.PacketInterface, callID uint32, gid types.UInt32, strMessage types.String, dontCareMyBlockList types.Bool, participationCount types.UInt16) (*nex.RMCMessage, *nex.Error) {
	globals.Logger.Infof("gid: %d, strMessage: %s, dontCareMyBlockList: %t, participationCount: %d\r\n", gid, strMessage, dontCareMyBlockList, participationCount)

	caller := types.NewPID(uint64(gid)) // Gathering ID and caller are the same here

	database.EndCallRinging(caller)

	sessionKey := types.NewBuffer(make([]byte, 32))

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	sessionKey.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.CallID = callID
	rmcResponse.MethodID = matchmake_extension.MethodJoinMatchmakeSessionEx

	return rmcResponse, nil
}
