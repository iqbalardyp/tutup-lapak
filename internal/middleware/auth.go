package custom_middleware

import (
	"net/http"

	customErrors "tutup-lapak/pkg/custom-errors"
	"tutup-lapak/pkg/dotenv"
	jwt "tutup-lapak/pkg/jwt"
	"tutup-lapak/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type AuthConfig struct {
	Env *dotenv.Env
}

func NewAuthMiddleware(env *dotenv.Env) *AuthConfig {
	return &AuthConfig{
		Env: env,
	}
}

func (a *AuthConfig) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			jwtToken, err := extractJWTTokenFromHeader(ctx.Request())
			if err != nil {
				return ctx.JSON(response.WriteErrorResponse(err))
			}

			claim, err := jwt.ClaimToken(jwtToken, a.Env.JWT_SECRET)
			if err != nil {
				err = errors.Wrap(customErrors.ErrUnauthorized, err.Error())
				return ctx.JSON(response.WriteErrorResponse(err))
			}

			ctx.Set("user", claim)

			// default user passing middleware if token is valid
			return next(ctx)
		}
	}
}

func extractJWTTokenFromHeader(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return "", errors.Wrap(customErrors.ErrUnauthorized, "missing auth token")
	}

	return authToken[len("Bearer "):], nil
}
