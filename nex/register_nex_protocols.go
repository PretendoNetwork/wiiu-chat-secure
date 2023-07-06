package nex

import (
	match_making "github.com/PretendoNetwork/nex-protocols-go/match-making"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/match-making-ext"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	nex_match_making "github.com/PretendoNetwork/wiiu-chat-secure/nex/match-making"
	nex_match_making_ext "github.com/PretendoNetwork/wiiu-chat-secure/nex/match-making-ext"
	nex_matchmake_extension "github.com/PretendoNetwork/wiiu-chat-secure/nex/matchmake-extension"
)

func registerNEXProtocols() {
	matchmakeExtensionProtocol := matchmake_extension.NewMatchmakeExtensionProtocol(globals.NEXServer)

	matchmakeExtensionProtocol.OpenParticipation(nex_matchmake_extension.OpenParticipation)
	matchmakeExtensionProtocol.CreateMatchmakeSession(nex_matchmake_extension.CreateMatchmakeSession)
	matchmakeExtensionProtocol.UpdateNotificationData(nex_matchmake_extension.UpdateNotificationData)
	matchmakeExtensionProtocol.GetFriendNotificationData(nex_matchmake_extension.GetFriendNotificationData)
	matchmakeExtensionProtocol.JoinMatchmakeSessionEx(nex_matchmake_extension.JoinMatchmakeSessionEx)

	matchMakingProtocol := match_making.NewMatchMakingProtocol(globals.NEXServer)

	matchMakingProtocol.UnregisterGathering(nex_match_making.UnregisterGathering)
	matchMakingProtocol.FindBySingleID(nex_match_making.FindBySingleID)
	matchMakingProtocol.GetSessionURLs(nex_match_making.GetSessionUrls)

	matchMakingExtProtocol := match_making_ext.NewMatchMakingExtProtocol(globals.NEXServer)

	matchMakingExtProtocol.EndParticipation(nex_match_making_ext.EndParticipation)
}
