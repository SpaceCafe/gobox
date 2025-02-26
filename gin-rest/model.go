package rest

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm"
)

var (
	_ types.IModel = (*Model)(nil)
)

// Model represents a base model with common fields for database records.
// To embed this base model into a model, consider following example:
//
//	    type Book struct {
//		       types.Model `yaml:",inline"`
//		       Title       string
//		       Author      string
//	    }
type Model struct {
	// ID is the primary key of the Model, represented as a UUID string.
	// It can only be written once upon creation.
	ID string `json:"id,omitempty" yaml:"id,omitempty" xml:"id,omitempty" gorm:"type:uuid;primary_key;<-:create" binding:"omitempty"`

	// CreatedAt is a timestamp indicating when the Model was created.
	CreatedAt *sql.NullTime `json:"created_at,omitempty" yaml:"created_at,omitempty" xml:"created_at,omitempty" gorm:"autoCreateTime" binding:"omitempty"`

	// UpdatedAt is a timestamp indicating when the Model was last updated.
	UpdatedAt *sql.NullTime `json:"updated_at,omitempty" yaml:"updated_at,omitempty" xml:"updated_at,omitempty" gorm:"autoUpdateTime" binding:"omitempty"`

	// DeletedAt is a timestamp indicating when the Model was deleted, used for soft deletes.
	DeletedAt *sql.NullTime `json:"deleted_at,omitempty" yaml:"deleted_at,omitempty" xml:"deleted_at,omitempty" gorm:"index" binding:"omitempty"`
}

// GetID returns the unique identifier of the model.
func (r *Model) GetID() string {
	return r.ID
}

// SetID sets the unique identifier for the model.
func (r *Model) SetID(id string) {
	r.ID = id
}

// BeforeCreate is a GORM hook that generates a UUID for the Model's ID field if it is not already set.
func (r *Model) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
