package nex

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go/v2"
	_ "github.com/PretendoNetwork/nex-protocols-go/v2"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func StartNEXServer() {
	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = true

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)
	globals.SecureEndpoint.DefaultStreamSettings.KeepAliveTimeout = 65535

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 4, 2))
	globals.SecureEndpoint.SetAccessKey("e7a47214")

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("==WiiU Chat - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("======================")
	})

	// * Register the common handlers first so that they can be overridden if needed
	registerCommonProtocols()
	registerNEXProtocols()

	globals.NEXServer.Listen(":" + os.Getenv("PN_WIIU_CHAT_SECURE_SERVER_PORT"))
}
