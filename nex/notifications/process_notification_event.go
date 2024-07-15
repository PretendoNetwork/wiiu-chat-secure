package nex_notifications

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	notifications "github.com/PretendoNetwork/nex-protocols-go/v2/notifications"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func ProcessNotificationEvent(callID uint32, pidSource uint32, uiType uint32, uiParam1 uint32, uiParam2 uint32, strParam string) {
	event := notifications_types.NewNotificationEvent()

	event.PIDSource = types.NewPID(uint64(pidSource)) // Sender PID
	event.Type.Value = uiType                         // Notification type
	event.Param1.Value = uiParam1                     // Gathering ID
	event.Param2.Value = uiParam2                     // Recipient PID
	event.StrParam.Value = strParam                   // Unknown

	eventObject := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureEndpoint.ByteStreamSettings()) // TODO: use the endpoint library versions/byte stream settings
	event.WriteTo(eventObject)

	rmcRequest := nex.NewRMCRequest(notifications.NewProtocol().Endpoint()) // TODO: check
	rmcRequestBytes := rmcRequest.Bytes()

	targetClient := globals.SecureEndpoint.FindConnectionByPID(uint64(uiParam2))

	requestPacket, _ := nex.NewPRUDPPacketV1(globals.SecureServer, targetClient, nil)

	requestPacket.SetType(constants.DataPacket)
	requestPacket.AddFlag(constants.PacketFlagNeedsAck)
	requestPacket.AddFlag(constants.PacketFlagReliable)
	requestPacket.SetPayload(rmcRequestBytes)

	globals.SecureServer.Send(requestPacket)
}
