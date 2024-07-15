package grpc

import (
	"context"
	"encoding/binary"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
	nintendo_notifications "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications"
	nintendo_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/nintendo-notifications/types"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"google.golang.org/grpc/metadata"
)

func SendFriendsNotification(caller uint32, target uint32, ringing bool) {
	ctx := metadata.NewOutgoingContext(context.Background(), globals.GRPCFriendsCommonMetadata)

	presence := friends_wiiu_types.NewNintendoPresenceV2()

	presence.Online.Value = true
	presence.GameKey = friends_wiiu_types.NewGameKey()
	presence.GameServerID.Value = 0x1005A000
	presence.PID = types.NewPID(1) // This is not a PID, but the amount of times the PID is repeated in bytes in the application data.
	presence.GatheringID.Value = caller

	if ringing {
		presence.Unknown2.Value = 0x65
	}

	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, target)

	presence.ApplicationData.Value = targetBytes

	presence.GameKey.TitleID.Value = 0x000500101005A100
	presence.GameKey.TitleVersion.Value = 55

	eventObject := nintendo_notifications_types.NewNintendoNotificationEvent()
	eventObject.Type.Value = nintendo_notifications.NotificationTypes.FriendStartedTitle
	eventObject.SenderPID = types.NewPID(uint64(caller))
	eventObject.DataHolder = types.NewAnyDataHolder()
	eventObject.DataHolder.TypeName = types.NewString("NintendoPresenceV2")
	eventObject.DataHolder.ObjectData = presence

	stream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())
	eventObject.WriteTo(stream)

	_, err := globals.GRPCFriendsClient.SendUserNotificationWiiU(ctx, &pb.SendUserNotificationWiiURequest{Pid: target, NotificationData: stream.Bytes()})
	if err != nil {
		globals.Logger.Criticalf("Greeting Friends gRPC failed! : %v", err)
	}
}
