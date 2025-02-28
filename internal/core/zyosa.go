package core

import (
	"fmt"
	"zyosa/internal/delivery/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Zyosa struct{
	Viper *viper.Viper
	App *fiber.App
	DB *gorm.DB
}

func Init(zyosa *Zyosa){
	route := routes.Route{
		App: zyosa.App,
	}

	route.Init()
}

func (a *Zyosa) Start(){
	err := a.App.Listen(fmt.Sprintf("%s:%s", a.Viper.GetString("app.host"), a.Viper.GetString("app.port")))
	if err != nil {
		logrus.Fatal(err)
	}
}