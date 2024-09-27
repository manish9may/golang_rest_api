package models

import (
	"database/sql"
	"fmt"
	"golang_rest_api/db"
	"log"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      uuid.UUID
}

func (e *Event) Save() error {
	log.Printf("Starting to save event: %+v\n", e)
	query := `INSERT INTO events(id,
	name, description, location, datetime, user_id)
	VALUES ($1, $2, $3, $4, $5,$6);`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatalf("Error Save database Step 1: %v", err)
		return err
	}
	e.ID = uuid.New()
	result, err := stmt.Exec(e.ID, e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		log.Fatalf("Error Save database Step 2: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected() // Check rows affected
	if err != nil {
		log.Fatalf("Error Save database Step 3: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were inserted")
	}
	return err
}

func GetAllEvents() ([]Event, error) {
	log.Println("Fetching all events from the database")
	query := `SELECT id, name, description, location, datetime, user_id FROM events`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("Error fetching events from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			log.Fatalf("Error scanning event from database: %v", err)
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error during row iteration: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched %d events", len(events))
	return events, nil
}

func FetchEventByID(id string) (*Event, error) {
	log.Printf("Fetching event with ID: %s", id)

	// Define the query to fetch the event by ID
	query := `SELECT id, name, description, location, datetime, user_id FROM events WHERE id = $1`

	// Execute the query
	row := db.DB.QueryRow(query, id)

	// Create an Event object to hold the result
	var event Event

	// Scan the result into the event struct
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No event found with ID: %s", id)
			return nil, nil // Return nil if no rows are found
		}
		log.Fatalf("Error fetching event from database: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched event: %+v", event)
	return &event, nil
}

// UpdateEvent updates an event by its ID in the database
func (e *Event) Update() error {
	log.Printf("Starting to update event: %+v\n", e)

	// Define the SQL query to update the event
	query := `
	UPDATE events 
	SET name = $1, description = $2, location = $3, datetime = $4, user_id = $5
	WHERE id = $6;
	`

	// Prepare the query statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatalf("Error updating database Step 1: %v", err)
		return err
	}

	// Execute the update with event data
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID, e.ID)
	if err != nil {
		log.Fatalf("Error updating database Step 2: %v", err)
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error updating database Step 3: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	log.Printf("Successfully updated event with ID: %s", e.ID)
	return nil
}

// Delete deletes an event by its ID from the database
func (e *Event) Delete() error {
	log.Printf("Starting to delete event with ID: %s\n", e.ID)

	// Define the SQL query to delete the event
	query := `DELETE FROM events WHERE id = $1;`

	// Prepare the query statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatalf("Error deleting from database Step 1: %v", err)
		return err
	}

	// Execute the delete statement with the event ID
	result, err := stmt.Exec(e.ID)
	if err != nil {
		log.Fatalf("Error deleting from database Step 2: %v", err)
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error deleting from database Step 3: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted, event not found")
	}

	log.Printf("Successfully deleted event with ID: %s", e.ID)
	return nil
}
