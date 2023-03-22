package demo

import (
	"mackey/pkg/config"
	"mackey/pkg/db"
	"mackey/pkg/logger"
	"mackey/pkg/tests"
)

type TestConfig struct {
	DB   db.Config
	Demo Config

	TestHash string `env:"TEST_HASH,defualt=test"`
}

func init() {
	tests.Provide(
		// configs
		config.Env(TestConfig{}),

		// logger
		logger.NewLogger,

		// db
		db.New,
		db.NewQueryHook,
		db.NewDemoRepo,

		New,
	)
}
