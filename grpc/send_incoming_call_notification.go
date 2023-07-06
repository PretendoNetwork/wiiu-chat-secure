package grpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/friends-wiiu/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications/types"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"google.golang.org/grpc/metadata"
)

func SendIncomingCallNotification(caller uint32, target uint32) {
	ctx := metadata.NewOutgoingContext(context.Background(), globals.GRPCFriendsCommonMetadata)

	presence := friends_wiiu_types.NewNintendoPresenceV2()

	presence.ChangedFlags = 0x1FF
	presence.Online = true
	presence.GameKey = friends_wiiu_types.NewGameKey()
	presence.Unknown2 = 0x65
	presence.GameServerID = 0x1005A000
	presence.PID = 4 // This is not a PID.
	presence.GatheringID = caller

	// The application data here is the PID, repeated 4 times. Why, you ask? Who knows!
	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, target)

	presence.ApplicationData = bytes.Repeat(targetBytes, 4)

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
		log.Fatalf("Greeting Friends gRPC failed! : %v", err)
	}
}
