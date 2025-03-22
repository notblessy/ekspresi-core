package main

import (
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
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

	cloudinary, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	continueOrFatal(err)

	userRepo := repository.NewUserRepository(postgres)
	uploaderRepo := repository.NewUploaderRepository(cloudinary, postgres)
	portfolioRepo := repository.NewPortfolioRepository(postgres, uploaderRepo)
	membershipRepo := repository.NewMembershipRepository(postgres)
	membershipPlanRepo := repository.NewMembershipPlanRepository(postgres)

	httpService := router.NewHTTPService()
	httpService.RegisterPostgres(postgres)
	httpService.RegisterUserRepository(userRepo)
	httpService.RegisterMembershipRepository(membershipRepo)
	httpService.RegisterMembershipPlanRepository(membershipPlanRepo)
	httpService.RegisterUploaderRepository(uploaderRepo)
	httpService.RegisterPortfolioRepository(portfolioRepo)

	httpService.Router(e)

	e.Logger.Fatal(e.Start(":3400"))
}

func continueOrFatal(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
