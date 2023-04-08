package nex_notifications

import (
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/notifications"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

// * This technically is not the right way to format an RMC handler
// * since it imples that the request format is "caller uint32, target uint32, callID uint32"
// * when really the data we are sending in this function is the
// * request format
// *
// * This is just a convenient way to handle sending notifications
// * since the server never gets requests for this protocol
// TODO - Maybe refactor this so that it's more in line with the spec?

func ProcessNotificationEvent(caller uint32, target uint32, callID uint32) {
	event := notifications.NewNotificationEvent()

	event.PIDSource = caller          // Sender PID
	event.Type = 101000               // Notification type
	event.Param1 = caller             // Gathering ID
	event.Param2 = target             // Recipient PID
	event.StrParam = "Invite Request" // Unknown

	eventObject := nex.NewStreamOut(globals.NEXServer)
	eventObject.WriteStructure(event)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(notifications.ProtocolID)
	rmcRequest.SetCallID(0xFFFF + callID)
	rmcRequest.SetMethodID(notifications.MethodProcessNotificationEvent)
	rmcRequest.SetParameters(eventObject.Bytes())

	rmcRequestBytes := rmcRequest.Bytes()

	targetClient := globals.NEXServer.FindClientFromPID(uint32(target))

	requestPacket, _ := nex.NewPacketV1(targetClient, nil)

	requestPacket.SetVersion(1)
	requestPacket.SetSource(0xA1)
	requestPacket.SetDestination(0xAF)
	requestPacket.SetType(nex.DataPacket)
	requestPacket.SetPayload(rmcRequestBytes)

	requestPacket.AddFlag(nex.FlagNeedsAck)
	requestPacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(requestPacket)
}
