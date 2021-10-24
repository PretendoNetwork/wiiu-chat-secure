package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server

var secureServer *nexproto.SecureProtocol

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(2)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetAccessKey("e7a47214")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==WiiU Chat - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("======================")
	})

	secureServer = nexproto.NewSecureProtocol(nexServer)
	matchmakeExtensionServer := nexproto.NewMatchmakeExtensionProtocol(nexServer)
	natTraversalServer := nexproto.NewNATTraversalProtocol(nexServer)

	matchmakeExtensionServer.OpenParticipation(openParticipation)
	matchmakeExtensionServer.CreateMatchmakeSession(createMatchmakeSession)
	matchmakeExtensionServer.UpdateNotificationData(updateNotificationData)
	matchmakeExtensionServer.GetFriendNotificationData(getFriendNotificationData)

	natTraversalServer.ReportNATProperties(reportNATProperties)

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	// Secure protocol handles

	secureServer.Register(register)

	nexServer.Listen(":60005")
}
