package main

import (
	"os"

	"github.com/ziadrahmatullah/minimarket-app/logger"
	"github.com/ziadrahmatullah/minimarket-app/migration"
	"github.com/ziadrahmatullah/minimarket-app/repository"
)

func main() {
	logger.SetLogrusLogger()

	_ = os.Setenv("APP_ENV", "debug")

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	migration.Migrate(db)
}
