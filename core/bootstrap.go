package core

import (
	"firebase-sso/db"
	"firebase-sso/helpers/env"
	"fmt"
)

func InitBackgroundServices() {
	// .ENV
	initEnv := env.Environment{EnvPath: ".env"}
	initEnv.LoadEnv()

	// DATABASE
	_, connErr := db.PostgresDB.SetupPostgresDB()
	if connErr != nil {
		fmt.Println(connErr.Error())
		return
	}

	// CHECK DB CONNECTION
	db.PostgresDB.PingDB()

}
