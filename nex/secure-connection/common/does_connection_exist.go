package nex_secure_connection_common

import (
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func DoesConnectionExist(rvcid uint32) bool {
	pid := globals.NEXServer.FindClientFromConnectionID(rvcid).PID()
	return database.DoesSessionExist(pid)
}
