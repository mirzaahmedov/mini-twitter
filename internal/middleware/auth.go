package middleware

import (
	"fmt"
	"net/http"
	"twitter/internal/utils"

	"github.com/labstack/echo/v4"
)

func Authentication(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				c.Logger().Error(err)
				c.NoContent(http.StatusUnauthorized)
				return err
			}
			if cookie.Value == "" {
				c.Logger().Error("no authorization cookie provided")
				c.NoContent(http.StatusUnauthorized)
				return fmt.Errorf("no authorization cookie provided")
			}

			claims, err := utils.VerifyToken(secretKey, cookie.Value)
			if err != nil {
				c.Logger().Error(err)
				return err
			}

			userID, ok := claims["user_id"]
			if !ok || userID == "" {
				c.Logger().Error("invalid authorization cookie provided")
				c.NoContent(http.StatusUnauthorized)
				return fmt.Errorf("invalid authorization cookie provided")
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}
