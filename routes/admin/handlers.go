package admin

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Starwish-Cottage/service/core"
)

func LoginHandler(c *gin.Context) {
	var request LoginRequest
	var VALID_DAYS = os.Getenv("LOGIN_DAYS")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := core.FirestoreClient
	doc, err := client.Collection("admin_users").Doc(request.Username).Get(c)
	if err != nil || doc.Data()["password"] != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to access Firestore"})
		return
	}

	parsedDays, err := strconv.ParseInt(VALID_DAYS, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid login days configuration"})
	}

	sessionToken, err := generateSessionToken(request.Username, time.Hour*time.Duration(24*parsedDays))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	fullName, _ := doc.Data()["full_name"].(string)

	resonse := LoginResponse{
		FullName:     fullName,
		SessionToken: sessionToken,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "data": resonse})
}

func VerifySession(c *gin.Context) {

}

func generateSessionToken(username string, expiration time.Duration) (string, error) {
	claim := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(JWT_SECRET)
}
