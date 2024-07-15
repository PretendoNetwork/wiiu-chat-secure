package nex

import (
	common_matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	nex_matchmake_extension "github.com/PretendoNetwork/wiiu-chat-secure/nex/matchmake-extension"
)

func registerNEXProtocols() {
	matchmakeExtensionProtocol := matchmake_extension.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	common_matchmake_extension.NewCommonProtocol(matchmakeExtensionProtocol)
	matchmakeExtensionProtocol.GetFriendNotificationData = nex_matchmake_extension.GetFriendNotificationData
	matchmakeExtensionProtocol.UpdateNotificationData = nex_matchmake_extension.UpdateNotificationData
}
