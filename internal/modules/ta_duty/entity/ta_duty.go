package entity

import "time"

type Holiday struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	ClassesID *int      `json:"classes_id"`
}
