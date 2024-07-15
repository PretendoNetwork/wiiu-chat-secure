package matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	notifications "github.com/PretendoNetwork/nex-protocols-go/v2/notifications"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
)

func GetFriendNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType *types.PrimitiveS32) (*nex.RMCMessage, *nex.Error) {
	dataList := types.NewList[*notifications_types.NotificationEvent]()

	caller, target, ringing := database.GetCallInfoByTarget(packet.Sender().PID().Value())

	// TODO: Multiple calls. Wii U Chat can handle it, but we don't support it yet
	if caller != 0 && target == packet.Sender().PID().Value() && ringing {
		// Being called
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.RequestJoinGathering, notifications.NotificationSubTypes.RequestJoinGathering.None)

		notification := notifications_types.NewNotificationEvent()

		notification.PIDSource = types.NewPID(caller)
		notification.Type.Value = notificationType
		notification.Param1.Value = uint32(caller)
		notification.Param2.Value = uint32(target)
		notification.StrParam.Value = "Invite Request"

		dataList.Append(notification)
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())
	dataList.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodGetFriendNotificationData
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
