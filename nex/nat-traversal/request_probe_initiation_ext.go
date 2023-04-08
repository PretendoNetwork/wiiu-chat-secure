package nex_nat_traversal

import (
	"fmt"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nat_traversal "github.com/PretendoNetwork/nex-protocols-go/nat-traversal"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func RequestProbeInitiationExt(err error, client *nex.Client, callID uint32, targetList []string, stationToProbe string) {
	fmt.Println(targetList)
	fmt.Println(stationToProbe)

	rmcResponse := nex.NewRMCResponse(nat_traversal.ProtocolID, callID)
	rmcResponse.SetSuccess(nat_traversal.MethodRequestProbeInitiationExt, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)

	// Send probe to other users

	rmcRequestStream := nex.NewStreamOut(globals.NEXServer)

	rmcRequestStream.WriteString(stationToProbe)

	rmcRequestBody := rmcRequestStream.Bytes()

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nat_traversal.ProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nat_traversal.MethodInitiateProbe)
	rmcRequest.SetParameters(rmcRequestBody)

	rmcRequestBytes := rmcRequest.Bytes()

	for _, target := range targetList {
		targetUrl := nex.NewStationURL(target)
		targetRVCID, _ := strconv.Atoi(targetUrl.RVCID())
		targetClient := globals.NEXServer.FindClientFromConnectionID(uint32(targetRVCID))
		fmt.Println(targetClient)
		if targetClient != nil {
			requestPacket, _ := nex.NewPacketV1(targetClient, nil)
			requestPacket.SetVersion(1)
			requestPacket.SetSource(0xA1)
			requestPacket.SetDestination(0xAF)
			requestPacket.SetType(nex.DataPacket)
			requestPacket.SetPayload(rmcRequestBytes)

			requestPacket.AddFlag(nex.FlagNeedsAck)
			requestPacket.AddFlag(nex.FlagReliable)

			globals.NEXServer.Send(requestPacket)
		}
	}
}
