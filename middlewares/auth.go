package middlewares

import (
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type ClaimsMap struct {
	jwt.StandardClaims
	Scope string
}

const SecretKey = "5144F11C-4A9E-4801-B005-66F1FB8CE027"

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsMap{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "unauthenticated",
		})
	}

	payload := token.Claims.(*ClaimsMap)
	isAmbassador := strings.Contains(c.Path(), "/api/v1/ambassadors")
	if (payload.Scope == "admin" && isAmbassador) || (payload.Scope == "ambassador" && !isAmbassador) {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	return c.Next()
}

func GenerateJWT(id uint, scope string) (string, error) {
	payload := ClaimsMap{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 1).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecretKey))

}

func GetUserID(c *fiber.Ctx) (uint, error) {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsMap{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	payload := token.Claims.(*ClaimsMap)
	id, _ := strconv.Atoi(payload.Subject)
	return uint(id), nil
}
