package nex

import (
	nat_traversal "github.com/PretendoNetwork/nex-protocols-common-go/nat-traversal"
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"

	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func registerCommonProtocols() {
	nat_traversal.NewCommonNATTraversalProtocol(globals.NEXServer)
	secureconnection.NewCommonSecureConnectionProtocol(globals.NEXServer)
}
