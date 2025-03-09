package main

import (
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/notblessy/ekspresi-core/db"
	"github.com/notblessy/ekspresi-core/repository"
	"github.com/notblessy/ekspresi-core/router"
	"github.com/notblessy/ekspresi-core/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("cannot load .env file")
	}

	postgres := db.NewPostgres()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Path",
		},
	}))
	e.Use(middleware.CORS())
	e.Validator = &utils.Ghost{Validator: validator.New()}

	userRepo := repository.NewUserRepository(postgres)

	httpService := router.NewHTTPService()
	httpService.RegisterPostgres(postgres)
	httpService.RegisterUserRepository(userRepo)

	httpService.Router(e)

	e.Logger.Fatal(e.Start(":3400"))
}
