package models

import (
	"errors"
	"time"

	"example.com/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `binding:"required" json:"name"`
	Description string    `binding:"required" json:"description"`
	Location    string    `binding:"required" json:"location"`
	DateTime    time.Time `binding:"required" json:"dateTime"`
	UserID      int64     `json:"userId"`
}

type Registration struct {
	ID   int64 `json:"id"`
	User struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
	EventID int64 `json:"eventId"`
}

func (e *Event) Save() error {
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetEventByID(id int64) (*Event, error) {
	query := `SELECT * FROM events WHERE id = ?`

	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)

	return err
}

func (e Event) Register(userId int64) error {
	query := `
	INSERT INTO registrations(event_id, user_id) VALUES (?, ?)
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e Event) GetAllRegistrations() ([]Registration, error) {
	registrations := []Registration{}

	query := `
	SELECT 
		r.id AS id,
		r.event_id AS event_id,
		u.id AS user_id,
		u.email AS user_email
	FROM 
		registrations r
	JOIN
		users u ON r.user_id = u.id
	WHERE event_id = ?
	`

	rows, err := db.DB.Query(query, e.ID)

	if err != nil {
		return registrations, err
	}

	defer rows.Close()

	for rows.Next() {
		var registration Registration

		err := rows.Scan(&registration.ID, &registration.EventID, &registration.User.ID, &registration.User.Email)

		if err != nil {
			return registrations, err
		}

		registrations = append(registrations, registration)
	}

	return registrations, nil
}

func (e Event) GetRegistrationByUser(userId int64) (Registration, error) {
	query := `
	SELECT 
		r.id AS id,
		r.event_id AS event_id,
		u.id AS user_id,
		u.email AS user_email
	FROM 
		registrations r
	JOIN
		users u ON r.user_id = u.id
	WHERE event_id = ? AND user_id = ?
	`

	row := db.DB.QueryRow(query, e.ID, userId)

	var registration Registration

	err := row.Scan(&registration.ID, &registration.EventID, &registration.User.ID, &registration.User.Email)

	if err != nil {
		return Registration{}, errors.New("registration not found")
	}

	return registration, nil
}

func (e Event) CancelRegistration(userId int64) error {
	query := `
	DELETE FROM registrations WHERE event_id = ? AND user_id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}
