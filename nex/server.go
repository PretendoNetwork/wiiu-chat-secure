package nex

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	_ "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func StartNEXServer() {
	globals.NEXServer = nex.NewServer()
	globals.NEXServer.SetPRUDPVersion(1)
	globals.NEXServer.SetPRUDPProtocolMinorVersion(2)
	globals.NEXServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	globals.NEXServer.SetAccessKey("e7a47214")
	globals.NEXServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 3,
		Minor: 4,
		Patch: 2,
	})
	globals.NEXServer.SetPingTimeout(65535)

	globals.NEXServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==WiiU Chat - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("======================")
	})

	globals.NEXServer.On("Kick", func(packet *nex.PacketV1) {
		fmt.Println("Kick client event called")
	})

	globals.NEXServer.On("Disconnect", func(packet *nex.PacketV1) {
		fmt.Println("Disconnect client event called")
	})

	// * Register the common handlers first so that they can be overridden if needed
	registerCommonProtocols()
	registerNEXProtocols()

	globals.NEXServer.Listen(":60005")
}
