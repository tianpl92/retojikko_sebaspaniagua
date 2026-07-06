package domain

import (
	"time"
)

// SavedProposal maps to the public_call_user_associations table.
type SavedProposal struct {
	ID              string     `json:"id"`
	PublicCallID    int        `json:"public_call_id"`
	UserID          string     `json:"user_id"`
	AssociationDate time.Time  `json:"association_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
