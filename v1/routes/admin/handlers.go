package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Starwish-Cottage/service/core"
	"github.com/Starwish-Cottage/service/v1/models"
)

func LoginHandler(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.LoginResponse{Success: false, FullName: "", SessionToken: "", Message: err.Error()})
		return
	}

	// check Firestore for the admin username and password
	client := core.FirestoreClient
	doc, err := client.Collection("admin_users").Doc(request.Username).Get(c)
	if err != nil || !doc.Exists() || doc.Data()["password"].(string) != request.Password {
		c.JSON(http.StatusUnauthorized, models.LoginResponse{Success: false, FullName: "", SessionToken: "", Message: "Incorrect username or password"})
		return
	}

	sessionToken, err := GenerateSessionToken(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.LoginResponse{Success: false, FullName: "", SessionToken: "", Message: "Failed to generate session token"})
		return
	}

	fullName, _ := doc.Data()["full_name"].(string)

	c.JSON(http.StatusOK, models.LoginResponse{Success: true, FullName: fullName, SessionToken: sessionToken, Message: "Login successful"})
}

func VerifySessionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, models.VerifySessionResponse{Valid: true, Message: "Session verification successful"})
}

func UploadImageHandler(c *gin.Context) {
	// add file type validation
	form, err := c.MultipartForm()

	username, _ := c.Get("username")

	if err != nil {
		c.JSON(http.StatusBadRequest, models.UploadImageResponse{Success: false, ImageUrls: nil, Message: err.Error()})
		return
	}
	files := form.File["files"]
	var urls []string
	for _, file := range files {
		destination := "./scripts/src_imgs/" + username.(string) + "/" + file.Filename
		if err := c.SaveUploadedFile(file, destination); err != nil {
			c.JSON(http.StatusInternalServerError, models.UploadImageResponse{Success: false, ImageUrls: nil, Message: err.Error()})
		}
		imageUrl := fmt.Sprintf("/images/%s/%s", username, file.Filename)
		urls = append(urls, imageUrl)
	}
	c.JSON(http.StatusOK, models.UploadImageResponse{Success: true, ImageUrls: urls, Message: "Image uploaded successfully"})
}
