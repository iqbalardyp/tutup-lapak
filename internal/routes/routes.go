package routes

import (
	"net/http"
	custom_middleware "tutup-lapak/internal/middleware"
	product_handler "tutup-lapak/internal/product/handler"
	"tutup-lapak/pkg/response"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App            *echo.Echo
	S3Uploader     *manager.Uploader
	Middleware     *custom_middleware.AuthConfig
	ProductHandler *product_handler.ProductHandler
}

func (r *RouteConfig) SetupRoutes() {
	r.App.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.BaseResponse{
			Status:  "Ok",
			Message: "",
		})
	})

	v1 := r.App.Group("/v1")
	r.setupPublicRoutes(v1)
	r.setupAuthRoutes(v1, r.Middleware.Authenticate())
}

func (r *RouteConfig) setupPublicRoutes(group *echo.Group) {
	group.GET("/product", r.ProductHandler.GetProducts)
}

func (r *RouteConfig) setupAuthRoutes(group *echo.Group, m echo.MiddlewareFunc) {
	r.setupProductAuthRoutes(group, m)
}

func (r *RouteConfig) setupProductAuthRoutes(group *echo.Group, m echo.MiddlewareFunc) {
	product := group.Group("/product")
	// product.POST("", r.ProductHandler.CreateProduct, m)
	// product.PATCH("/:productId", r.ProductHandler.UpdateProduct, m)
	// product.DELETE("/:productId", r.ProductHandler.DeleteProduct, m)
	product.POST("", r.ProductHandler.CreateProduct)
	product.PATCH("/:productId", r.ProductHandler.UpdateProduct)
	product.DELETE("/:productId", r.ProductHandler.DeleteProduct)
}
