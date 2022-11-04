package v1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type UserModel struct {
	BaseModel
	Username string `gorm:"uniqueIndex"`
}

type SessionModel struct {
	BaseModel
	UserID uuid.UUID
	//User      UserModel
	StartTime time.Time
	EndTime   time.Time
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}