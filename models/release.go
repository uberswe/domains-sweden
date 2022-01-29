package models

import (
	"gorm.io/gorm"
	"time"
)

type Release struct {
	gorm.Model
	ReleasedAt *time.Time `gorm:"uniqueIndex"`
	Domains    []Domain   `gorm:"many2many:domain_releases;"`
}
