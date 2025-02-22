package nex_matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/wiiu-chat/database"
	"github.com/PretendoNetwork/wiiu-chat/globals"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	notifications "github.com/PretendoNetwork/nex-protocols-go/v2/notifications"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
)

func GetFriendNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType types.Int32) (*nex.RMCMessage, *nex.Error) {
	dataList := types.NewList[notifications_types.NotificationEvent]()

	caller, target, ringing := database.GetCallInfoByTarget(packet.Sender().PID())

	// TODO: Multiple calls. Wii U Chat can handle it, but we don't support it yet
	if caller != 0 && target == packet.Sender().PID() && ringing {
		// Being called
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.RequestJoinGathering, notifications.NotificationSubTypes.RequestJoinGathering.None)

		notification := notifications_types.NewNotificationEvent()

		notification.PIDSource = caller
		notification.Type = types.NewUInt32(notificationType)
		notification.Param1 = types.NewUInt32(uint32(caller))
		notification.Param2 = types.NewUInt32(uint32(target))
		notification.StrParam = "Invite Request"

		dataList = append(dataList, notification)
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	dataList.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.CallID = callID
	rmcResponse.MethodID = matchmake_extension.MethodGetFriendNotificationData

	return rmcResponse, nil
}
