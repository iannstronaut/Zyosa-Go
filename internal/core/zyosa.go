package core

import (
	"fmt"
	"zyosa/internal/delivery/middleware"
	"zyosa/internal/delivery/routes"
	"zyosa/internal/domains/user"
	"zyosa/internal/domains/user/repository"
	"zyosa/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Zyosa struct {
	Viper      *viper.Viper
	App        *fiber.App
	DB         *gorm.DB
	JWTService *services.JWTService
}

func Init(zyosa *Zyosa) {
	jwtService := services.NewJWTService(zyosa.Viper.GetString("app.secret"))
	userRepo := repository.NewUserRepository(zyosa.DB)
	userHandler := user.NewHandler(userRepo, zyosa.Viper, jwtService)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	route := routes.Route{
		App:            zyosa.App,
		UserRoute:      userHandler,
		JWTService:     jwtService,
		AuthMiddleware: authMiddleware,
	}

	route.Init()
}

func (a *Zyosa) Start() {
	err := a.App.Listen(fmt.Sprintf("%s:%s", a.Viper.GetString("app.host"), a.Viper.GetString("app.port")))
	if err != nil {
		logrus.Fatal(err)
	}
}
