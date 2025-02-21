package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/PretendoNetwork/wiiu-chat/globals"
)

var Postgres *sql.DB

func ConnectPostgres() {
	var err error

	Postgres, err = sql.Open("postgres", os.Getenv("PN_WUC_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	globals.Logger.Success("Connected to Postgres!")

	InitPostgres()
}
