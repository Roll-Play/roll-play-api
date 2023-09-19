package middlewares

import (
	"errors"
	"net/http"
	"strings"

	api_error "github.com/Roll-play/roll-play-backend/pkg/errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, api_error.Error{
				Error:   "Missing Authorization header",
				Message: http.StatusText(http.StatusUnauthorized),
			})
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secretKey := []byte("your-secret-key")

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return secretKey, nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, api_error.Error{
				Error:   "Invalid token",
				Message: http.StatusText(http.StatusUnauthorized),
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user, ok := claims["sub"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, api_error.Error{
					Error:   "Invalid token",
					Message: http.StatusText(http.StatusUnauthorized),
				})
			}

			c.Set("user", user)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, api_error.Error{
			Error:   "Invalid token",
			Message: http.StatusText(http.StatusUnauthorized),
		})
	}
}
