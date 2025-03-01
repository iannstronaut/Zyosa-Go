package core

import (
	"fmt"
	"zyosa/internal/delivery/middleware"
	"zyosa/internal/delivery/routes"
	UserHandler"zyosa/internal/domains/user"
	UserRepository "zyosa/internal/domains/user/repository"
	AdminRepository "zyosa/internal/domains/admin/repository"
	AdminHandler "zyosa/internal/domains/admin"
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

	userRepo := UserRepository.NewUserRepository(zyosa.DB)
	userHandler := UserHandler.NewHandler(userRepo, zyosa.Viper, jwtService)

	adminRepo := AdminRepository.NewAdminRepository(zyosa.DB)
	adminHandler := AdminHandler.NewHandler(adminRepo, zyosa.Viper, jwtService)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	route := routes.Route{
		App:            zyosa.App,
		UserRoute:      userHandler,
		AdminRoute:     adminHandler,
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
