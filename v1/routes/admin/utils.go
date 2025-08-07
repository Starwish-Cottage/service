package admin

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
func GenerateSessionToken(username string, expiration time.Duration) (string, error) {
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
func GetValidHours() time.Duration {
	daysStr := os.Getenv("LOGIN_DAYS")
	days, err := strconv.ParseInt(daysStr, 10, 64)
	if err != nil || days <= 0 {
		days = 7 // Default to 7 days
	}
	return time.Hour * time.Duration(days*24)
}
