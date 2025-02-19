package nex

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	common_match_making "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	common_match_making_ext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	common_matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	common_nat_traversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	mm_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nat_traversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/PretendoNetwork/wiiu-chat/database"
	local_globals "github.com/PretendoNetwork/wiiu-chat/globals"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
)


func updateNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType types.UInt32, uiParam1 types.UInt32, uiParam2 types.UInt32, strParam types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}
	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodUpdateNotificationData
	rmcResponse.CallID = callID
	return rmcResponse, nil
}
func getFriendNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType types.Int32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	dataList := types.NewList[*notifications_types.NotificationEvent]()

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	dataList.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodGetFriendNotificationData
	rmcResponse.CallID = callID

	return rmcResponse, nil
}

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	secure := common_secure.NewCommonProtocol(secureProtocol)
	secure.CreateReportDBRecord = func(pid types.PID, reportID types.UInt32, reportData types.QBuffer) error {
		return nil
	}

	matchmakingManager := common_globals.NewMatchmakingManager(local_globals.SecureEndpoint, database.Postgres)

	natTraversalProtocol := nat_traversal.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	common_nat_traversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := match_making.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonMatchMakingProtocol := common_match_making.NewCommonProtocol(matchMakingProtocol)
	commonMatchMakingProtocol.SetManager(matchmakingManager)
	

	matchMakingExtProtocol := match_making_ext.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol := common_match_making_ext.NewCommonProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol.SetManager(matchmakingManager)
	

	matchmakeExtensionProtocol := matchmake_extension.NewProtocol()
	local_globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol := common_matchmake_extension.NewCommonProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol.SetManager(matchmakingManager)
	matchmakeExtensionProtocol.SetHandlerGetFriendNotificationData(getFriendNotificationData)
	matchmakeExtensionProtocol.SetHandlerUpdateNotificationData(updateNotificationData)

	commonMatchmakeExtensionProtocol.CleanupSearchMatchmakeSession = func(matchmakeSession *mm_types.MatchmakeSession) {}
	commonMatchmakeExtensionProtocol.OnAfterAutoMatchmakeWithSearchCriteriaPostpone = func(packet nex.PacketInterface, lstSearchCriteria types.List[mm_types.MatchmakeSessionSearchCriteria], anyGathering types.AnyObjectHolder[mm_types.GatheringInterface], strMessage types.String) {
	}
	commonMatchmakeExtensionProtocol.CleanupMatchmakeSessionSearchCriterias = func(searchCriterias types.List[mm_types.MatchmakeSessionSearchCriteria]) {}
}
