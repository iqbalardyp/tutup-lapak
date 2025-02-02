package config

import (
	"time"
	"tutup-lapak/db"
	file_handler "tutup-lapak/internal/file/handler"
	file_repository "tutup-lapak/internal/file/repository"
	file_usecase "tutup-lapak/internal/file/usecase"
	custom_middleware "tutup-lapak/internal/middleware"
	product_handler "tutup-lapak/internal/product/handler"
	product_repository "tutup-lapak/internal/product/repository"
	product_usecase "tutup-lapak/internal/product/usecase"
	purchase_handler "tutup-lapak/internal/purchase/handler"
	purchase_repository "tutup-lapak/internal/purchase/repository"
	purchase_usecase "tutup-lapak/internal/purchase/usecase"
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

	productRepo := product_repository.NewProductRepo(config.DB.Pool)
	productUsecase := product_usecase.NewProductUsecase(productRepo)
	productHandler := product_handler.NewProductHandler(productUsecase, config.Validator)

	purchaseRepo := purchase_repository.NewPurchaseRepository(config.DB.Pool)
	purchaseUsecase := purchase_usecase.NewPurchaseUseCase(purchaseRepo, productRepo)
	purchaseHandler := purchase_handler.NewPurchaseHandler(purchaseUsecase, config.Validator)

	fileRepo := file_repository.NewFileRepository(config.DB.Pool)
	fileUsecase := file_usecase.NewFileUseCase(config.S3Uploader, config.Env, fileRepo)
	fileHandler := file_handler.NewFileHandler(fileUsecase, config.Log)

	routes := routes.RouteConfig{
		App:             config.App,
		S3Uploader:      config.S3Uploader,
		Middleware:      authMiddleware,
		ProductHandler:  productHandler,
		PurchaseHandler: purchaseHandler,
		FileHandler:     fileHandler,
	}

	routes.SetupRoutes()
}
