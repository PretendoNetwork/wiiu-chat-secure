package grpc

import (
	"context"
	"encoding/binary"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"google.golang.org/grpc/metadata"
)

func SendFriendsNotification(caller uint32, target uint32, ringing bool) {
	ctx := metadata.NewOutgoingContext(context.Background(), globals.GRPCFriendsCommonMetadata)

	presence := friends_wiiu_types.NewNintendoPresenceV2()

	presence.ChangedFlags = 0x1FF
	presence.Online = true
	presence.GameKey = friends_wiiu_types.NewGameKey()
	presence.GameServerID = 0x1005A000
	presence.PID = 1 // This is not a PID, but the amount of times the PID is repeated in bytes in the application data.
	presence.GatheringID = caller

	if ringing {
		presence.Unknown2 = 0x65
	}

	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, target)

	presence.ApplicationData = targetBytes

	presence.GameKey.TitleID = 0x000500101005A100
	presence.GameKey.TitleVersion = 55

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type = nintendo_notifications.NotificationTypes.FriendStartedTitle
	eventObject.SenderPID = caller
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoPresenceV2")
	eventObject.DataHolder.SetObjectData(presence)

	stream := nex.NewStreamOut(globals.NEXServer)
	eventObjectBytes := eventObject.Bytes(stream)

	_, err := globals.GRPCFriendsClient.SendUserNotificationWiiU(ctx, &pb.SendUserNotificationWiiURequest{Pid: target, NotificationData: eventObjectBytes})
	if err != nil {
		globals.Logger.Criticalf("Greeting Friends gRPC failed! : %v", err)
	}
}
