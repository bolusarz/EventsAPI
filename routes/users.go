package routes

import (
	"net/http"

	"example.com/models"
	"example.com/utils"
	"github.com/gin-gonic/gin"
)

func signup(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user data."})
		return
	}

	err = user.Save()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not save user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": user})
}

func login(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user data."})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "data": token})
}
