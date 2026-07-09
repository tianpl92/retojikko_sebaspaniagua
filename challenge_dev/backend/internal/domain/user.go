package domain

import "time"

type User struct {
	ID        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Gender    string     `json:"gender"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone_number"`
	Password  string     `json:"-"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
