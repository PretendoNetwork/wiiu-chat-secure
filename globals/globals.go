package globals

import (
	pb_accounts "github.com/PretendoNetwork/grpc-go/account"
	pb_friends "github.com/PretendoNetwork/grpc-go/friends"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/plogger-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Logger *plogger.Logger
var KerberosPassword = "password" // * Default password
var AuthenticationServer *nex.PRUDPServer
var AuthenticationEndpoint *nex.PRUDPEndPoint
var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint
var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb_accounts.AccountClient
var GRPCAccountCommonMetadata metadata.MD
var GRPCFriendsClientConnection *grpc.ClientConn
var GRPCFriendsClient pb_friends.FriendsClient
var GRPCFriendsCommonMetadata metadata.MD
