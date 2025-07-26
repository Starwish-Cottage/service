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

// LoginHandler handles the admin user login process.
// It expects a JSON payload containing the username and password.
// The function authenticates the user against Firestore, and if successful,
// generates a session token and returns the user's full name and session token in the response.
// On authentication failure or error, it returns an appropriate HTTP error response.
//
// Request JSON:
//
//	{
//	  "username": "admin",
//	  "password": "password"
//	}
//
// Response JSON (on success):
//
//	{
//	   "full_name": "Admin User",
//	   "session_token": "token_string"
//	}
//
// Possible HTTP status codes:
//   - 200 OK: Login successful
//   - 400 Bad Request: Invalid request payload
//   - 401 Unauthorized: Authentication failed
//   - 500 Internal Server Error: Session token generation failed
func LoginHandler(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check Firestore for the admin username and password
	client := core.FirestoreClient
	doc, err := client.Collection("admin_users").Doc(request.Username).Get(c)
	if err != nil || !doc.Exists() || doc.Data()["password"].(string) != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	sessionToken, err := generateSessionToken(request.Username, getValidHours())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	fullName, _ := doc.Data()["full_name"].(string)

	c.JSON(http.StatusOK, gin.H{"full_name": fullName, "session_token": sessionToken})
}

// VerifySessionHandler handles the verification of a user's session token.
// It expects a JSON payload containing a session token, validates the token's
// signature and expiration, and returns an appropriate HTTP response.
// If the token is missing, invalid, or expired, it responds with an error.
// On successful validation, it responds with a success message.
//
// Expected request body:
//
//	{
//	  "sessionToken": "<JWT token string>"
//	}
//
// Responses:
//
//	200 OK    - Session validated successfully
//	400 Bad Request - Missing or malformed session token
//	401 Unauthorized - Invalid or expired session token
func VerifySessionHandler(c *gin.Context) {
	var request VerifySessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenStr := request.SessionToken
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session token is required"})
		return
	}

	var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		}

		// Check if the session token is still valid
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session token has expired"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session validated successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
	}
}

// generateSessionToken generates a JWT session token for the given username with the specified expiration duration.
// It uses the secret key from the environment variable "JWT_SECRET" to sign the token.
//
// Parameters:
//   - username: The username for which the token is generated.
//   - expiration: The duration for which the token is valid.
//
// Returns:
//   - A signed JWT token string.
//   - An error if token generation fails.
func generateSessionToken(username string, expiration time.Duration) (string, error) {
	claim := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(JWT_SECRET)
}

// getValidHours retrieves the session token validity duration in hours.
// It reads the "LOGIN_DAYS" environment variable to determine the number of valid days.
// If the variable is not set or invalid, it defaults to 7 days.
//
// Returns:
//   - A time.Duration representing the validity period in hours.
func getValidHours() time.Duration {
	daysStr := os.Getenv("LOGIN_DAYS")
	days, err := strconv.ParseInt(daysStr, 10, 64)
	if err != nil || days <= 0 {
		days = 7 // Default to 7 days
	}
	return time.Hour * time.Duration(days*24)
}
