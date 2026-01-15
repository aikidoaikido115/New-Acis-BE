package middlewares

import (
    "fmt"

    "github.com/aikidoaikido115/New-Acis-BE/configs"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(config configs.JWT) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        tokenString := ctx.Get("Authorization")
        if tokenString == "" || len(tokenString) < 8 {
            return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "status":  "Unauthorized",
                "message": "Missing or invalid authorization header",
            })
        }

        tokenString = tokenString[7:] // Remove "Bearer "

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // ✅ STRICT CHECK: Only accept HS256
            if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            
            return []byte(config.Secret), nil
        })

        if err != nil || !token.Valid {
            return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "status":  "Unauthorized",
                "message": "Invalid or expired token",
            })
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            ctx.Locals("user_id", claims["user_id"])
        }

        return ctx.Next()
    }
}