package models

import (
	"gorm.io/gorm"
	"time"
)

type Domain struct {
	gorm.Model
	Host         string `gorm:"index:idx_domains_host;size:256;type:varchar(256)"`
	RegisteredAt *time.Time
	ExpiresAt    *time.Time
	Nameservers  []Nameserver `gorm:"many2many:domain_nameservers;"`
	Releases     []Release    `gorm:"many2many:domain_releases;"`
	Parses       []Parse
}
