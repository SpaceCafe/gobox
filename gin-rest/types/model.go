package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model represents a base model with common fields for database records.
// To embed this base model into a model, consider following example:
//
//	    type BookModel struct {
//		       types.Model `yaml:",inline"`
//		       Title       string
//		       Author      string
//	    }
type Model struct {

	// ID is the primary key of the Model, represented as a UUID string.
	// It can only be written once upon creation.
	ID string `json:"id" yaml:"id" xml:"id" gorm:"type:uuid;primary_key;<-:create"`

	// CreatedAt is a timestamp indicating when the Model was created.
	CreatedAt time.Time `json:"created_at" yaml:"created_at" xml:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt is a timestamp indicating when the Model was last updated.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at" xml:"updated_at" gorm:"autoUpdateTime"`

	// DeletedAt is a timestamp indicating when the Model was deleted, used for soft deletes.
	DeletedAt *sql.NullTime `json:"deleted_at" yaml:"deleted_at" xml:"deleted_at" gorm:"index"`
}

// BeforeCreate is a GORM hook that generates a UUID for the Model's ID field if it is not already set.
func (r *Model) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
