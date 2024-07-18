package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/PretendoNetwork/grpc-go/friends"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var logger = plogger.NewLogger()

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Warning("Error loading .env file")
	}

	kerberosPassword := os.Getenv("PN_WIIU_CHAT_KERBEROS_PASSWORD")
	friendsGRPCHost := os.Getenv("PN_WIIU_CHAT_FRIENDS_GRPC_HOST")
	friendsGRPCPort := os.Getenv("PN_WIIU_CHAT_FRIENDS_GRPC_PORT")
	friendsGRPCAPIKey := os.Getenv("PN_WIIU_CHAT_FRIENDS_GRPC_API_KEY")

	if strings.TrimSpace(kerberosPassword) == "" {
		globals.Logger.Warningf("PN_KERBEROS_PASSWORD environment variable not set. Using default password: %q", globals.KerberosPassword)
	} else {
		globals.KerberosPassword = kerberosPassword
	}

	if strings.TrimSpace(friendsGRPCHost) == "" {
		globals.Logger.Error("PN_WIIU_CHAT_FRIENDS_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(friendsGRPCPort) == "" {
		globals.Logger.Error("PN_WIIU_CHAT_FRIENDS_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(friendsGRPCPort); err != nil {
		globals.Logger.Errorf("PN_WIIU_CHAT_FRIENDS_GRPC_PORT is not a valid port. Expected 0-65535, got %s", friendsGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_WIIU_CHAT_FRIENDS_GRPC_PORT is not a valid port. Expected 0-65535, got %s", friendsGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(friendsGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_WIIU_CHAT_FRIENDS_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCFriendsClientConnection, err = grpc.Dial(fmt.Sprintf("%s:%s", friendsGRPCHost, friendsGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to friends gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCFriendsClient = pb.NewFriendsClient(globals.GRPCFriendsClientConnection)
	globals.GRPCFriendsCommonMetadata = metadata.Pairs(
		"X-API-Key", friendsGRPCAPIKey,
	)

	globals.InitAccounts()
	database.ConnectAll()
}
