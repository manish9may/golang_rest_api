package middleware

import (
	"fmt"
	"golang_rest_api/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTAuthMiddleware checks if the JWT token is valid
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
			ctx.Abort()
			return
		}

		// The token is prefixed with "Bearer ", so we strip that
		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return models.JwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid JWT token: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(string)
			ctx.Set("userID", userID) // Store user ID in context
		}

		ctx.Next()
	}
}
