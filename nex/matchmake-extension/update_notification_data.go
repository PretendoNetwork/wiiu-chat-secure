package nex_matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"github.com/PretendoNetwork/wiiu-chat-secure/grpc"
	nex_notifications "github.com/PretendoNetwork/wiiu-chat-secure/nex/notifications"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/nex-protocols-go/notifications"
)

func UpdateNotificationData(err error, client *nex.Client, callID uint32, uiType uint32, uiParam1 uint32, uiParam2 uint32, strParam string) {
	globals.Logger.Infof("uiType: %d, uiParam1: %d, uiParam2: %d, strParam: %s\r\n", uiType, uiParam1, uiParam2, strParam)
	recipientClient := globals.NEXServer.FindClientFromPID(uiParam2)

	if uiType == notifications.NotificationCategories.RequestJoinGathering {
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.RequestJoinGathering, notifications.NotificationSubTypes.RequestJoinGathering.None)
		database.NewCall(client.PID(), uiParam2)

		// If they don't have a session with the app, tell Friends to alert them on the HOME menu.
		if recipientClient != nil && recipientClient.StationURLs() != nil {
			nex_notifications.ProcessNotificationEvent(callID, client.PID(), notificationType, uiParam1, uiParam2, strParam)
		} else {
			grpc.SendIncomingCallNotification(client.PID(), uiParam2)
		}
	}

	if uiType == notifications.NotificationCategories.EndGathering {
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.EndGathering, notifications.NotificationSubTypes.EndGathering.None)
		database.EndCall(uiParam1)

		// Alert the other side we aren't calling anymore.

		if recipientClient != nil && recipientClient.StationURLs() != nil {
			nex_notifications.ProcessNotificationEvent(callID, client.PID(), notificationType, uiParam1, uiParam2, strParam)
		}

		// The user must be kicked, otherwise the app hangs forever.
		globals.NEXServer.TimeoutKick(client)
	}

	// Build response packet
	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodUpdateNotificationData, nil)

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
