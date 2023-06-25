package grpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"os"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/friends/wiiu"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/nintendo-notifications"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func SendIncomingCallNotification(caller uint32, target uint32) {
	// Connect to Friends gRPC.
	conn, err := grpc.Dial(os.Getenv("FRIENDS_GRPC_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection to Friends gRPC failed! : %v", err)
	}

	defer conn.Close()
	c := pb.NewFriendsClient(conn)

	md := metadata.Pairs(
		"X-API-Key", os.Getenv("FRIENDS_GRPC_API_KEY"),
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	presence := friends_wiiu.NewNintendoPresenceV2()

	presence.ChangedFlags = 0x1FF
	presence.Online = true
	presence.GameKey = friends_wiiu.NewGameKey()
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

	eventObject := nintendo_notifications.NewNintendoNotificationEvent()
	eventObject.Type = nintendo_notifications.NotificationTypes.FriendStartedTitle
	eventObject.SenderPID = caller
	eventObject.DataHolder = nex.NewDataHolder()
	eventObject.DataHolder.SetTypeName("NintendoPresenceV2")
	eventObject.DataHolder.SetObjectData(presence)

	stream := nex.NewStreamOut(globals.NEXServer)
	eventObjectBytes := eventObject.Bytes(stream)

	_, err = c.SendUserNotificationWiiU(ctx, &pb.SendUserNotificationWiiURequest{Pid: target, NotificationData: eventObjectBytes})
	if err != nil {
		log.Fatalf("Greeting Friends gRPC failed! : %v", err)
	}
}
