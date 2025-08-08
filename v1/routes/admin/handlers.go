package admin

import (
	"fmt"
	"net/http"
	"os"
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
//		 "message": "Login successful / failed reason"
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
		c.JSON(http.StatusUnprocessableEntity, LoginResponse{false, "", "", err.Error()})
		return
	}

	// check Firestore for the admin username and password
	client := core.FirestoreClient
	doc, err := client.Collection("admin_users").Doc(request.Username).Get(c)
	if err != nil || !doc.Exists() || doc.Data()["password"].(string) != request.Password {
		c.JSON(http.StatusUnauthorized, LoginResponse{false, "", "", "Incorrect username or password"})
		return
	}

	sessionToken, err := GenerateSessionToken(request.Username, GetValidHours())
	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{false, "", "", "Failed to generate session token"})
		return
	}

	fullName, _ := doc.Data()["full_name"].(string)

	c.JSON(http.StatusOK, LoginResponse{true, fullName, sessionToken, "Login successful"})
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
//	{
//	   "valid": false
//	   "message": "session token expired"
//	}
//
//	200 OK    - Session validated successfully
//	400 Bad Request - Missing or malformed session token
//	401 Unauthorized - Invalid or expired session token
func VerifySessionHandler(c *gin.Context) {
	var request VerifySessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	tokenStr := request.SessionToken

	var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, VerifySessionResponse{false, err.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, VerifySessionResponse{false, "Invalid session token"})
			return
		}

		// Check if the session token is still valid
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, VerifySessionResponse{false, "Session token expired, please sign in again"})
			return
		}

		c.JSON(http.StatusOK, VerifySessionResponse{true, "Session verification successful"})
	} else {
		c.JSON(http.StatusUnauthorized, VerifySessionResponse{false, "Invalid session token"})
	}
}

// UploadImageHandler handles the upload of a single image file.
// It accepts a multipart form request with an image file under the "image" field,
// saves the file to the local filesystem, and returns the accessible URL for the uploaded image.
//
// Request:
//   - Method: POST
//   - Content-Type: multipart/form-data
//   - Form field: "image" (file)
//
// Request Example:
//   curl -X POST http://localhost:8080/v1/admin/upload-image \
//     -F "image=@/path/to/image.jpg"
//
// Response JSON (on success):
//
//	{
//	  "image_url": "/images/filename.jpg",
//	  "message": "Image uploaded successfully"
//	}
//
// Response JSON (on error):
//
//	{
//	  "image_url": "",
//	  "message": "error description"
//	}
//
// Possible HTTP status codes:
//   - 200 OK: Image uploaded successfully
//   - 400 Bad Request: No file provided or invalid form data
//   - 500 Internal Server Error: Failed to save file to filesystem
//
// File Storage:
//   - Files are saved to: ./scripts/src_imgs/
//   - Accessible via URL: /images/{filename}
//   - Note: Requires static file serving configuration for /images route
//
// Security Considerations:
//   - TODO: Add file type validation (currently accepts any file type)
//   - TODO: Add file size limits
//   - TODO: Sanitize filename to prevent path traversal attacks
//   - TODO: Add authentication/authorization checks
//
// Known Issues:
//   - Missing return statement after SaveUploadedFile error (function continues execution)
//   - No file type validation implemented
//   - Potential filename conflicts (no unique naming strategy)
//   - Directory traversal vulnerability if filename contains "../"

func UploadImageHandler(c *gin.Context) {
	// add file type validation
	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, UploadImageResponse{false, nil, err.Error()})
		return
	}
	files := form.File["files"]
	var urls []string
	for _, file := range files {
		destination := "./scripts/src_imgs/" + file.Filename
		if err := c.SaveUploadedFile(file, destination); err != nil {
			c.JSON(http.StatusInternalServerError, UploadImageResponse{false, nil, err.Error()})
		}
		imageUrl := fmt.Sprintf("/images/%s", file.Filename)
		urls = append(urls, imageUrl)
	}
	c.JSON(http.StatusOK, UploadImageResponse{true, urls, "Image uploaded successfully"})
}
