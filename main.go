package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eymen-iron/web-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var SECRET = []byte("super-secret-auth-key")
var api_key = "api-key"

func CreateJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenStr, err := token.SignedString(SECRET)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func ValidateJWT(c *fiber.Ctx) error {

	tokenString := c.Get("Token")

	if tokenString != "" {
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, c.JSON(fiber.Map{"success": false, "message": "Invalid token"})
			}
			return SECRET, nil
		})

		if err != nil {
			return c.JSON(fiber.Map{"success": false, "message": "Invalid token"})
		}

		if token.Valid {
			return c.Next()
		}
	} else {
		return c.JSON(fiber.Map{"success": false, "message": "Token is required"})
	}

	return nil
}

func GetJwt(c *fiber.Ctx) error {
	if access := c.Get("Access"); access != "" {
		if access != api_key {
			return c.JSON(fiber.Map{"success": false, "message": "Invalid access key"})
		}
		token, err := CreateJWT()
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"success": true, "message": "Invalid access key", "token": token})
	}

	return nil
}

func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "message": "Welcome to the home page"})
}

func main() {

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: " http://localhost:3000, http://localhost:3001",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	api := app.Group("/api")

	// api

	// home
	api.Use(ValidateJWT)
	api.Get("/all-posts", router.GetAllPosts)
	api.Get("/category/:name", router.GetPostsByCategory)
	api.Get("/post/:name", router.GetPostByName)

	app.Get("/jwt", GetJwt)

	app.Listen(":3500")
}
