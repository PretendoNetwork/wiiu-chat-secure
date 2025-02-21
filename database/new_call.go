package database

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func NewCall(caller types.PID, target types.PID) {
	_, err := Postgres.Exec(`INSERT INTO ongoingcalls (caller_pid, target_pid, ringing) VALUES ($1, $2, $3);`, caller, target, true)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}
}
