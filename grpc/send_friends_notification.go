package grpc

import (
	"context"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
	"github.com/PretendoNetwork/wiiu-chat/globals"
	"google.golang.org/grpc/metadata"
)

func SendFriendsNotification(caller types.PID, target types.PID, ringing types.Bool) {
	ctx := metadata.NewOutgoingContext(context.Background(), globals.GRPCFriendsCommonMetadata)

	presence := friends_wiiu_types.NewNintendoPresenceV2()

	presence.Online = true
	presence.GameKey = friends_wiiu_types.NewGameKey()
	presence.GameServerID = 0x1005A000
	presence.PID = 1 // This is not a PID, but the amount of times the PID is repeated in bytes in the application data.
	presence.GatheringID = types.NewUInt32(uint32(caller))

	if ringing {
		presence.Unknown2 = types.NewUInt32(0x65)
	}

	targetBytes := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	target.WriteTo(targetBytes)

	presence.ApplicationData = targetBytes.Bytes()

	presence.GameKey.TitleID = 0x000500101005A100
	presence.GameKey.TitleVersion = 55

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.SenderPID = caller
	eventObject.DataHolder = types.NewDataHolder()
	eventObject.DataHolder.Object = presence

	stream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)
	eventObject.WriteTo(stream)
	eventObjectBytes := stream.Bytes()

	_, err := globals.GRPCFriendsClient.SendUserNotificationWiiU(ctx, &pb.SendUserNotificationWiiURequest{Pid: uint32(target), NotificationData: eventObjectBytes})
	if err != nil {
		globals.Logger.Criticalf("Greeting Friends gRPC failed! : %v", err)
	}
}
