package main

import (
	"mackey/pkg/app"
)

func main() {
	app.Build().Run()

	//err := app.Invoke(func(manager *demo.Manager, db db.DB, logger *logger.Logger) error {
	//	ver, err := db.Version()
	//	if err != nil {
	//		return err
	//	}
	//
	//	logger.Info(context.Background(), fmt.Sprintf("running on %s version", ver))
	//
	//	return manager.Run()
	//})
}
