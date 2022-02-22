package models

import (
	"gorm.io/gorm"
	"time"
)

type ParseEvent struct {
	gorm.Model
	ParseID uint
	Parse   Parse
	URL     string
	Type    string
	Time    time.Time
}
