package app

import (
	"context"
	"fmt"

	"mackey/pkg/config"
	"mackey/pkg/db"
	"mackey/pkg/demo"
	"mackey/pkg/logger"

	"go.uber.org/dig"
)

type Config struct {
	DB   db.Config
	Demo demo.Config
}

type App struct {
	*dig.Container
}

func Build() *App {
	container := dig.New()

	// config
	container.Provide(config.Env(Config{}))

	// logger
	container.Provide(logger.NewLogger)

	// db
	container.Provide(db.New)
	container.Provide(db.NewQueryHook)
	container.Provide(db.NewDemoRepo)

	// managers
	container.Provide(demo.New)

	// server etc

	return &App{Container: container}
}

func (a *App) Run() {
	err := a.Invoke(func(manager *demo.Manager, db db.DB, logger *logger.Logger) error {
		ver, err := db.Version()
		if err != nil {
			return err
		}

		logger.Info(context.Background(), fmt.Sprintf("running on %s version", ver))

		return manager.Run()
	})

	if err != nil {
		panic(err)
	}
}
