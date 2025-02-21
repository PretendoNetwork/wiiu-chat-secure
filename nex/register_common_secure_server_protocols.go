package nex

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_nat_traversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nat_traversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	local_globals "github.com/PretendoNetwork/wiiu-chat/globals"
	local_match_making "github.com/PretendoNetwork/wiiu-chat/nex/match-making"
	local_match_making_ext "github.com/PretendoNetwork/wiiu-chat/nex/match-making-ext"
	local_matchmake_extension "github.com/PretendoNetwork/wiiu-chat/nex/matchmake-extension"
)

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	secure := common_secure.NewCommonProtocol(secureProtocol)
	secure.CreateReportDBRecord = func(pid types.PID, reportID types.UInt32, reportData types.QBuffer) error {
		return nil
	}

	natTraversalProtocol := nat_traversal.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	common_nat_traversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := match_making.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	matchMakingProtocol.FindBySingleID = local_match_making.FindBySingleID
	matchMakingProtocol.GetSessionURLs = local_match_making.GetSessionUrls
	matchMakingProtocol.UnregisterGathering = local_match_making.UnregisterGathering

	matchMakingExtProtocol := match_making_ext.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	matchMakingExtProtocol.EndParticipation = local_match_making_ext.EndParticipation

	matchmakeExtensionProtocol := matchmake_extension.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	matchmakeExtensionProtocol.GetFriendNotificationData = local_matchmake_extension.GetFriendNotificationData
	matchmakeExtensionProtocol.UpdateNotificationData = local_matchmake_extension.UpdateNotificationData
	matchmakeExtensionProtocol.CreateMatchmakeSession = local_matchmake_extension.CreateMatchmakeSession
	matchmakeExtensionProtocol.JoinMatchmakeSessionEx = local_matchmake_extension.JoinMatchmakeSessionEx
	matchmakeExtensionProtocol.OpenParticipation = local_matchmake_extension.OpenParticipation
}
