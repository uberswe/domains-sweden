package models

import (
	"gorm.io/gorm"
	"time"
)

type Parse struct {
	gorm.Model
	Content           []byte  `gorm:"type:longblob"`
	Screenshot        []byte  `gorm:"type:blob"`
	BlurredScreenshot []byte  `gorm:"type:blob"` // There is a lot of porn out there
	Size              float64 // in Mb
	LoadTime          float64 // in Seconds
	Error             *string
	Requested         time.Time
	Events            []ParseEvent
	DomainID          uint
	Domain            Domain
}
