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

	rmcMessage := nex.RMCRequest{}
	rmcMessage.SetProtocolID(nexproto.NATTraversalProtocolID)
	rmcMessage.SetCallID(0xffff0000 + callID)
	rmcMessage.SetMethodID(nexproto.NATTraversalMethodInitiateProbe)
	rmcRequestStream := nex.NewStreamOut(nexServer)
	rmcRequestStream.WriteString(stationToProbe)
	rmcRequestBody := rmcRequestStream.Bytes()
	rmcMessage.SetParameters(rmcRequestBody)
	rmcMessageBytes := rmcMessage.Bytes()

	for _, target := range targetList {
		targetUrl := nex.NewStationURL(target)
		targetPID, _ := strconv.Atoi(targetUrl.PID())
		targetClient := nexServer.FindClientFromPID(uint32(targetPID))
		fmt.Println(targetClient)
		if targetClient != nil {
			messagePacket, _ := nex.NewPacketV1(targetClient, nil)
			messagePacket.SetVersion(1)
			messagePacket.SetSource(0xA1)
			messagePacket.SetDestination(0xAF)
			messagePacket.SetType(nex.DataPacket)
			messagePacket.SetPayload(rmcMessageBytes)

			messagePacket.AddFlag(nex.FlagNeedsAck)
			messagePacket.AddFlag(nex.FlagReliable)

			nexServer.Send(messagePacket)
		}
	}
}
