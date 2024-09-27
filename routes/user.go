package routes

import (
	"golang_rest_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createUser(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload!"})
		return
	}
	err = user.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.ID})
}

func loginUser(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind the incoming JSON to loginData struct
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data"})
		return
	}

	// Initialize a User model
	user := models.User{}

	// Call the Login method, which returns a JWT if successful
	token, err := user.Login(loginData.Email, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Respond with the JWT token
	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
