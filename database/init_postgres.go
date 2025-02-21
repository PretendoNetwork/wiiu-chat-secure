package database

import "github.com/PretendoNetwork/wiiu-chat/globals"

func InitPostgres() {
	var err error

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS ongoingcalls (
		caller_pid integer UNIQUE PRIMARY KEY,
		target_pid integer,
		ringing bool
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	// Empty out call list if there happens to be old data lingering
	_, err = Postgres.Exec(`DELETE FROM ongoingcalls;`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return
	}

	globals.Logger.Success("Postgres tables created")
}
