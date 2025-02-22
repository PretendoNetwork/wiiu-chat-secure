package nex_notifications

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/constants"
	"github.com/PretendoNetwork/nex-go/v2/types"
	notifications "github.com/PretendoNetwork/nex-protocols-go/v2/notifications"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func ProcessNotificationEvent(callID uint32, packet nex.PacketInterface, uiType types.UInt32, uiParam1 types.UInt32, uiParam2 types.UInt32, strParam types.String) {
	event := notifications_types.NewNotificationEvent()

	event.PIDSource = packet.Sender().PID() // Sender PID
	event.Type = uiType                     // Notification type
	event.Param1 = uiParam1                 // Gathering ID
	event.Param2 = uiParam2                 // Recipient PID
	event.StrParam = strParam               // Unknown

	eventObject := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	event.WriteTo(eventObject)

	rmcRequest := nex.NewRMCRequest(globals.SecureEndpoint)
	rmcRequest.ProtocolID = notifications.ProtocolID
	rmcRequest.CallID = 0xFFFF + callID
	rmcRequest.MethodID = notifications.MethodProcessNotificationEvent
	rmcRequest.Parameters = eventObject.Bytes()

	rmcRequestBytes := rmcRequest.Bytes()

	targetClient := globals.SecureEndpoint.FindConnectionByPID(uint64(uiParam2))

	var responsePacket nex.PacketInterface
	var prudpPacket nex.PRUDPPacketInterface

	endpoint := globals.SecureEndpoint
	server := endpoint.Server
	prudpPacket, _ = nex.NewPRUDPPacketV1(server, targetClient, nil)

	prudpPacket.SetType(constants.DataPacket)
	prudpPacket.AddFlag(constants.PacketFlagReliable)
	prudpPacket.AddFlag(constants.PacketFlagNeedsAck)
	prudpPacket.SetSourceVirtualPortStreamType(packet.(*nex.PRUDPPacketV1).DestinationVirtualPortStreamType())
	prudpPacket.SetSourceVirtualPortStreamID(packet.(*nex.PRUDPPacketV1).DestinationVirtualPortStreamID())
	prudpPacket.SetDestinationVirtualPortStreamType(packet.(*nex.PRUDPPacketV1).SourceVirtualPortStreamType())
	prudpPacket.SetDestinationVirtualPortStreamID(packet.(*nex.PRUDPPacketV1).SourceVirtualPortStreamID())
	prudpPacket.SetSubstreamID(packet.(*nex.PRUDPPacketV1).SubstreamID())

	responsePacket = prudpPacket
	responsePacket.SetPayload(rmcRequestBytes)

	packet.Sender().Endpoint().Send(responsePacket)
}
