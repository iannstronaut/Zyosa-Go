package main

import (
	_ "zyosa/docs"
	"zyosa/internal/config"
	"zyosa/internal/core"

	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

var (
	ZyosaApp *core.Zyosa
)

// @title Zyosa API
// @version 1.0
// @description Zyosa Api with Go Fiber
// @host localhost:5000
// @BasePath /api
func main() {
	viper := config.NewConfig()

	app := core.NewFiber(viper)

	app.Get("/docs/*", swagger.HandlerDefault)

	db, err := core.NewDB(viper)
	if err != nil {
		logrus.Fatalf("unable to initialize database: %s", err.Error())
	}

	ZyosaApp = &core.Zyosa{
		Viper: viper,
		App:   app,
		DB:    db,
	}

	core.Init(ZyosaApp)

	ZyosaApp.Start()
}
