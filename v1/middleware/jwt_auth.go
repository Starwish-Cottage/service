package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Starwish-Cottage/service/v1/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ProcessJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Authorization header required"})
		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Bearer token required"})
		c.Abort()
		return
	}

	var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Invalid token - failed to parse token"})
		c.Abort()
		return
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Token expired, please signin again"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Missing token field - exp"})
			c.Abort()
			return
		}

		// Extract the username and set in context
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		} else {
			c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Required field not found in token - username"})
			c.Abort()
			return
		}

	} else {
		c.JSON(http.StatusUnauthorized, models.VerifySessionResponse{Valid: false, Message: "Invalid token - failed to map claim"})
		c.Abort()
		return
	}

	// continue to next handler
	c.Next()
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return ProcessJWT
}
