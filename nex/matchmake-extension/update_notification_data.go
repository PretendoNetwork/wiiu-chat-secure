package nex_matchmake_extension

import (
	"log"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"github.com/PretendoNetwork/wiiu-chat-secure/grpc"
	nex_notifications "github.com/PretendoNetwork/wiiu-chat-secure/nex/notifications"

	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
)

func UpdateNotificationData(err error, client *nex.Client, callID uint32, uiType uint32, uiParam1 uint32, uiParam2 uint32, strParam string) {
	log.Printf("uiType: %d, uiParam1: %d, uiParam2: %d, strParam: %s\r\n", uiType, uiParam1, uiParam2, strParam)

	if uiType == 101 {
		database.NewCall(client.PID(), uiParam2)

		// If they don't have a session with the app, tell Friends to alert them on the HOME menu.
		if database.DoesSessionExist(uiParam2) {
			nex_notifications.ProcessNotificationEvent(callID, client.PID(), uiType, uiParam1, uiParam2, strParam)
		} else {
			grpc.SendIncomingCallNotification(client.PID(), uiParam2)
		}
	}

	if uiType == 102 {
		database.EndCall(uiParam1)

		// Alert the other side we aren't calling anymore.
		if database.DoesSessionExist(uiParam2) {
			nex_notifications.ProcessNotificationEvent(callID, client.PID(), uiType, uiParam1, uiParam2, strParam)
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
