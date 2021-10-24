package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func reportNATProperties(err error, client *nex.Client, callID uint32, natmapping uint32, natfiltering uint32, rtt uint32) {
	// TODO: Implement this

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.NATTraversalProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.NATTraversalMethodReportNATProperties, nil)

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
}
