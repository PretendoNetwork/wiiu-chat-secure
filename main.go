package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"

	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetPRUDPProtocolMinorVersion(2)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	nexServer.SetAccessKey("e7a47214")
	nexServer.SetNexVersion(30000)
	nexServer.SetPingTimeout(65535)

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==WiiU Chat - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("======================")
	})

	nexServer.On("Kick", func(packet *nex.PacketV1) {
		fmt.Println("Kick client event called")
		deletePlayerSession(packet.Sender().PID())
	})

	nexServer.On("Disconnect", func(packet *nex.PacketV1) {
		fmt.Println("Disconnect client event called")
		deletePlayerSession(packet.Sender().PID())
	})

	secureServer := secureconnection.NewCommonSecureConnectionProtocol(nexServer)
	matchmakeExtensionServer := nexproto.NewMatchmakeExtensionProtocol(nexServer)
	natTraversalServer := nexproto.NewNATTraversalProtocol(nexServer)
	matchMakingServer := nexproto.NewMatchMakingProtocol(nexServer)

	matchmakeExtensionServer.OpenParticipation(openParticipation)
	matchmakeExtensionServer.CreateMatchmakeSession(createMatchmakeSession)
	matchmakeExtensionServer.UpdateNotificationData(updateNotificationData)
	matchmakeExtensionServer.GetFriendNotificationData(getFriendNotificationData)
	matchmakeExtensionServer.JoinMatchmakeSessionEx(joinMatchmakeSessionEx)

	natTraversalServer.ReportNATProperties(reportNATProperties)
	natTraversalServer.RequestProbeInitiationExt(requestProbeInitiationExt)
	natTraversalServer.ReportNATTraversalResult(reportNATTraversalResult)

	matchMakingServer.UnregisterGathering(unregisterGathering)
	matchMakingServer.FindBySingleID(findBySingleID)
	matchMakingServer.GetSessionURLs(getSessionUrls)

	secureServer.AddConnection(func(rvcid uint32, urls []string, ip, port string) {
		pid := nexServer.FindClientFromConnectionID(rvcid).PID()
		addPlayerSession(pid, urls, ip, port)
	})
	secureServer.UpdateConnection(func(rvcid uint32, urls []string, ip, port string) {
		pid := nexServer.FindClientFromConnectionID(rvcid).PID()
		updatePlayerSessionAll(pid, urls, ip, port)
	})
	secureServer.DoesConnectionExist(func(rvcid uint32) bool {
		pid := nexServer.FindClientFromConnectionID(rvcid).PID()
		return doesSessionExist(pid)
	})

	nexServer.Listen(":60005")
}
