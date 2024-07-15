package nex

import (
	common_matchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	common_matchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	common_nattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func registerCommonProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	common_secure.NewCommonProtocol(secureProtocol)

	natTraversalProtocol := nattraversal.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	common_nattraversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	common_matchmaking.NewCommonProtocol(matchMakingProtocol)

	matchMakingExtProtocol := matchmakingext.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	common_matchmakingext.NewCommonProtocol(matchMakingExtProtocol)
}
