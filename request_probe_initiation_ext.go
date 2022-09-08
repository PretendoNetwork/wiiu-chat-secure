package main

import (
	"fmt"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"

	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func requestProbeInitiationExt(err error, client *nex.Client, callID uint32, targetList []string, stationToProbe string) {
	fmt.Println(targetList)
	fmt.Println(stationToProbe)

	rmcResponse := nex.NewRMCResponse(nexproto.NATTraversalProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.NATTraversalMethodRequestProbeInitiationExt, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)

	// Send probe to other users

	rmcRequestStream := nex.NewStreamOut(nexServer)

	rmcRequestStream.WriteString(stationToProbe)

	rmcRequestBody := rmcRequestStream.Bytes()

	rmcRequest := nex.NewRMCRequest()
	rmcRequest.SetProtocolID(nexproto.NATTraversalProtocolID)
	rmcRequest.SetCallID(3810693103)
	rmcRequest.SetMethodID(nexproto.NATTraversalMethodInitiateProbe)
	rmcRequest.SetParameters(rmcRequestBody)

	rmcRequestBytes := rmcRequest.Bytes()

	for _, target := range targetList {
		targetUrl := nex.NewStationURL(target)
		targetRVCID, _ := strconv.Atoi(targetUrl.RVCID())
		targetClient := nexServer.FindClientFromConnectionID(uint32(targetRVCID))
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

			nexServer.Send(requestPacket)
		}
	}
}
