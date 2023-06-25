package nex_notifications

import (
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/notifications"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func ProcessNotificationEvent(callID uint32, pidSource uint32, uiType uint32, uiParam1 uint32, uiParam2 uint32, strParam string) {
	event := notifications.NewNotificationEvent()

	event.PIDSource = pidSource // Sender PID
	event.Type = uiType * 1000  // Notification type, multiplied 1000 times because... yes?
	event.Param1 = uiParam1     // Gathering ID
	event.Param2 = uiParam2     // Recipient PID
	event.StrParam = strParam   // Unknown

	eventObject := nex.NewStreamOut(globals.NEXServer)
	eventObject.WriteStructure(event)

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(notifications.ProtocolID)
	rmcRequest.SetCallID(0xFFFF + callID)
	rmcRequest.SetMethodID(notifications.MethodProcessNotificationEvent)
	rmcRequest.SetParameters(eventObject.Bytes())

	rmcRequestBytes := rmcRequest.Bytes()

	targetClient := globals.NEXServer.FindClientFromPID(uiParam2)

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
