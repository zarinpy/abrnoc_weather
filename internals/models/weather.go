package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Weather struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CityName    string    `gorm:"not null" json:"cityName" validate:"required"`
	Country     string    `gorm:"not null" json:"country" validate:"required"`
	Temperature float64   `json:"temperature" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Humidity    int       `json:"humidity" validate:"required,min=0,max=100"`
	WindSpeed   float64   `json:"windSpeed" validate:"required,gte=0"`
	FetchedAt   time.Time `json:"fetchedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (w *Weather) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New()
	return nil
}
