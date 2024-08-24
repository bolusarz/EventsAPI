package routes

import (
	"net/http"
	"strconv"

	"example.com/models"
	"github.com/gin-gonic/gin"
)

func getRegistrations(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event := models.Event{}
	event.ID = eventId

	registrations, err := event.GetAllRegistrations()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch event"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Fetched successfully", "data": registrations})
}

func registerForEvent(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch event"})
		return
	}

	if event.UserID == userId {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You can not register for your event"})
		return
	}

	_, err = event.GetRegistrationByUser(userId)

	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You have registered already"})
		return
	}

	err = event.Register(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Registered!"})
}

func cancelRegistration(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not cancel registration."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Cancelled!"})
}
