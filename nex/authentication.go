package nex

import (
	"fmt"
	"os"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

var serverBuildString string

func StartAuthenticationServer() {
	serverBuildString = "branch:ngs_2_30 build:2_22_11148_30_2"

	globals.AuthenticationServer = nex.NewPRUDPServer()

	globals.AuthenticationEndpoint = nex.NewPRUDPEndPoint(1)
	globals.AuthenticationEndpoint.ServerAccount = globals.AuthenticationServerAccount
	globals.AuthenticationEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.AuthenticationEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.AuthenticationServer.BindPRUDPEndPoint(globals.AuthenticationEndpoint)
	globals.AuthenticationServer.ByteStreamSettings.UseStructureHeader = false

	// Technically this title is 3.4.2; however, it uses older-style structures and is therefore defined here as 3.3.2
	globals.AuthenticationServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 3, 2))
	globals.AuthenticationServer.AccessKey = "e7a47214"

	globals.AuthenticationEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== WUC - Auth ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("==================")
	})

	registerCommonAuthenticationServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_WUC_AUTHENTICATION_SERVER_PORT"))

	globals.AuthenticationServer.Listen(port)
}
