package config

import (
	"time"
	"tutup-lapak/db"
	custom_middleware "tutup-lapak/internal/middleware"
	"tutup-lapak/internal/routes"
	"tutup-lapak/pkg/dotenv"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	Env        *dotenv.Env
	App        *echo.Echo
	DB         *db.Postgres
	Log        *logrus.Logger
	Validator  *validator.Validate
	S3Uploader *manager.Uploader
}

func Bootstrap(config *BootstrapConfig) {
	// * Middleware
	config.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timeout",
		Timeout:      30 * time.Second,
	}))

	authMiddleware := custom_middleware.NewAuthMiddleware(config.Env)
	routes := routes.RouteConfig{
		App:        config.App,
		S3Uploader: config.S3Uploader,
		Middleware: authMiddleware,
	}

	routes.SetupRoutes()
}
