package middleware

import (
	"fmt"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")

		if token == "" {
			return errors.NewError(http.StatusUnauthorized, "unAuthorized token not present in the header")
		}
		fmt.Println(token)
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expires, ok := claims["expires"].(float64)
		if !ok {
			return errors.ErrUnAuthorized()
		}

		if int64(expires) < time.Now().Unix() {
			return errors.NewError(http.StatusUnauthorized, "token expired")
		}
		userID, ok := claims["id"].(string)
		if !ok {
			return errors.ErrUnAuthorized()
		}
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return errors.ErrUnAuthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println()
			return nil, errors.NewError(http.StatusUnauthorized, "invalid signing method "+fmt.Sprint(token.Header["alg"]))

		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, errors.NewError(404, "JWT_SECRET is not set")
		}
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println()
		return nil, errors.NewError(http.StatusUnauthorized, "failed to parse JWT token:"+fmt.Sprint(err))
	}
	if !token.Valid {

		return nil, errors.NewError(http.StatusUnauthorized, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrUnAuthorized()
	}
	return claims, nil
}
