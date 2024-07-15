package matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"github.com/PretendoNetwork/wiiu-chat-secure/grpc"
	nex_notifications "github.com/PretendoNetwork/wiiu-chat-secure/nex/notifications"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	notifications "github.com/PretendoNetwork/nex-protocols-go/v2/notifications"
)

func UpdateNotificationData(err error, packet nex.PacketInterface, callID uint32, uiType *types.PrimitiveU32, uiParam1 *types.PrimitiveU32, uiParam2 *types.PrimitiveU32, strParam *types.String) (*nex.RMCMessage, *nex.Error) {
	globals.Logger.Infof("uiType: %d, uiParam1: %d, uiParam2: %d, strParam: %s\r\n", uiType, uiParam1, uiParam2, strParam)
	recipientClient := globals.SecureEndpoint.FindConnectionByPID(uint64(uiParam2.Value))

	if uiType.Value == notifications.NotificationCategories.RequestJoinGathering {
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.RequestJoinGathering, notifications.NotificationSubTypes.RequestJoinGathering.None)
		database.NewCall(uint32(packet.Sender().PID().Value()), uiParam2.Value)

		// If they don't have a session with the app, tell Friends to alert them on the HOME menu.
		if recipientClient != nil && recipientClient.StationURLs != nil {
			nex_notifications.ProcessNotificationEvent(callID, uint32(packet.Sender().PID().Value()), notificationType, uiParam1.Value, uiParam2.Value, strParam.Value)
		} else {
			grpc.SendFriendsNotification(uint32(packet.Sender().PID().Value()), uiParam2.Value, true)
		}
	}

	if uiType.Value == notifications.NotificationCategories.EndGathering {
		notificationType := notifications.BuildNotificationType(notifications.NotificationCategories.EndGathering, notifications.NotificationSubTypes.EndGathering.None)
		database.EndCall(uiParam1.Value)

		// Alert the other side we aren't calling anymore.

		if recipientClient != nil && recipientClient.StationURLs != nil {
			nex_notifications.ProcessNotificationEvent(callID, uint32(packet.Sender().PID().Value()), notificationType, uiParam1.Value, uiParam2.Value, strParam.Value)
		} else {
			grpc.SendFriendsNotification(uint32(packet.Sender().PID().Value()), uiParam2.Value, false)
		}

		// The user must be kicked, otherwise the app hangs forever.
		// globals.NEXServer.TimeoutKick(client)
		// TODO: Use timeout manager
	}

	// Build response packet

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodUpdateNotificationData
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
