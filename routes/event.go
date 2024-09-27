package routes

import (
	"golang_rest_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func fetchEvents(ctx *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch events", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"events": events})
}

func createEvent(ctx *gin.Context) {
	var event models.Event
	err := ctx.ShouldBindJSON(&event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data!"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string)) // Ensure it's in UUID format
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid User ID format"})
		return
	}

	// Set the user ID for the event
	event.UserID = userUUID // Assign the UUID

	event.Save()
	ctx.JSON(http.StatusCreated, gin.H{"message": "Event Created", "event": event})
}

func fetchEvent(ctx *gin.Context) {
	// Get the event ID from the URL parameters
	eventID := ctx.Param("id")

	// Fetch the event by ID
	event, err := models.FetchEventByID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch event", "error": err.Error()})
		return
	}

	// If no event is found, respond with 404
	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	// Respond with the event
	ctx.JSON(http.StatusOK, gin.H{"event": event})
}

func updateEvent(ctx *gin.Context) {
	// Get the event ID from the URL parameters
	eventID := ctx.Param("id")

	// Fetch the event by ID (to ensure it exists before updating)
	event, err := models.FetchEventByID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch event", "error": err.Error()})
		return
	}

	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	// Bind the incoming JSON data to the event object
	err = ctx.ShouldBindJSON(event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the data!"})
		return
	}

	// Update the event
	err = event.Update()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update event", "error": err.Error()})
		return
	}

	// Respond with the updated event
	ctx.JSON(http.StatusOK, gin.H{"message": "Event updated successfully", "event": event})
}

func deleteEvent(ctx *gin.Context) {
	// Get the event ID from the URL parameters
	eventID := ctx.Param("id")

	// Fetch the event by ID (to ensure it exists before deleting)
	event, err := models.FetchEventByID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch event", "error": err.Error()})
		return
	}

	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	// Delete the event
	err = event.Delete()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete event", "error": err.Error()})
		return
	}

	// Respond with a success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
