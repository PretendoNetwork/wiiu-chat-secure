package nex_secure_connection_common

import (
	"github.com/PretendoNetwork/wiiu-chat-secure/database"
	"github.com/PretendoNetwork/wiiu-chat-secure/globals"
)

func UpdateConnection(rvcid uint32, urls []string, ip, port string) {
	pid := globals.NEXServer.FindClientFromConnectionID(rvcid).PID()
	database.UpdatePlayerSessionAll(pid, urls, ip, port)
}
