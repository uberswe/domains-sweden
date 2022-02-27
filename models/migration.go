package models

import (
	"gorm.io/gorm"
)

type Migration struct {
	gorm.Model
	Key string `gorm:"type:text"`
}
