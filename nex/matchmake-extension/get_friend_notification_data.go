package nex_matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/nex-protocols-go/notifications"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/notifications/types"
)

func GetFriendNotificationData(err error, client *nex.Client, callID uint32, uiType int32) {
	dataList := make([]*notifications_types.NotificationEvent, 0)

	caller, target, ringing := database.GetCallInfoByTarget(client.PID())

	// TODO: Multiple calls. Wii U Chat can handle it, but we don't support it yet
	if caller != 0 && target == client.PID() && ringing {
		// Being called
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.RequestJoinGathering, notifications.NotificationSubTypes.RequestJoinGathering.None)

		notification := notifications_types.NewNotificationEvent()

		notification.PIDSource = caller
		notification.Type = notificationType
		notification.Param1 = caller
		notification.Param2 = target
		notification.StrParam = "Invite Request"

		dataList = append(dataList, notification)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	rmcResponseStream.WriteListStructure(dataList)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodGetFriendNotificationData, rmcResponseBody)

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
