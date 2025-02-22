package database

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func EndCallRinging(caller types.PID) {
	_, err := Postgres.Exec(`UPDATE ongoingcalls SET ringing = $1 WHERE caller_pid = $2;`, false, caller)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}
}
