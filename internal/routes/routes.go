package routes

import (
	"net/http"
	custom_middleware "tutup-lapak/internal/middleware"
	"tutup-lapak/pkg/response"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App        *echo.Echo
	S3Uploader *manager.Uploader
	Middleware *custom_middleware.AuthConfig
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
}

func (r *RouteConfig) setupAuthRoutes(group *echo.Group, m echo.MiddlewareFunc) {
}
