package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

// NotificationEvent is a NotificationEvent
type NotificationEvent struct {
	sourcePID       uint32
	typeParameter   uint32
	parameter1      uint64
	parameter2      uint64
	stringParameter string
	nex.Structure
}

// Bytes encodes the NotificationEvent and returns a byte array
func (notificationEvent *NotificationEvent) Bytes(stream *nex.StreamOut) []byte {
	stream.WriteUInt32LE(notificationEvent.sourcePID)
	stream.WriteUInt32LE(notificationEvent.typeParameter)
	stream.WriteUInt64LE(notificationEvent.parameter1)
	stream.WriteUInt64LE(notificationEvent.parameter2)
	stream.WriteString(notificationEvent.stringParameter)

	return stream.Bytes()
}

func getFriendNotificationData(err error, client *nex.Client, callID uint32, uiType int32) {

	rmcResponseStream := nex.NewStreamOut(nexServer)

	// List<NotificationEvent>
	notificationEvent := &NotificationEvent{}

	notificationEvent.sourcePID = 1743126339                  // Sender PID
	notificationEvent.typeParameter = 102000                  // Notification type
	notificationEvent.parameter1 = 1                          // Gathering ID
	notificationEvent.parameter2 = 1730592963                 // Recipient PID
	notificationEvent.stringParameter = "Invite Cancellation" // Unknown

	notificationEvent2 := &NotificationEvent{}

	notificationEvent2.sourcePID = 1743126339             // Sender PID
	notificationEvent2.typeParameter = 101000             // Notification type
	notificationEvent2.parameter1 = 1                     // Gathering ID
	notificationEvent2.parameter2 = 1730592963            // Recipient PID
	notificationEvent2.stringParameter = "Invite Request" // Unknown

	rmcResponseStream.WriteUInt32LE(2)
	rmcResponseStream.WriteStructure(notificationEvent)
	rmcResponseStream.WriteStructure(notificationEvent2)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodGetFriendNotificationData, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
