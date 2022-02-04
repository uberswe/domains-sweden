package models

import (
	"gorm.io/gorm"
)

type Sitemap struct {
	gorm.Model
	Identifier string `gorm:"uniqueIndex;size:128"`
	Content    string `gorm:"type:longtext"`
}
