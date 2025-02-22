package database

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/wiiu-chat/globals"
)

func GetCallInfoByTarget(target types.PID) (caller_pid types.PID, target_pid types.PID, ringing types.Bool) {
	row := Postgres.QueryRow(`SELECT (caller_pid, target_pid, ringing) FROM ongoingcalls WHERE target_pid = $1;`, target)
	err := row.Scan(&caller_pid, &target_pid, &ringing)
	if err != nil {
		globals.Logger.Warning(err.Error())
		return 0, 0, false
	}
	return caller_pid, target_pid, ringing
}
