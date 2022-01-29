package models

import (
	"gorm.io/gorm"
)

type Nameserver struct {
	gorm.Model
	Host    string   `gorm:"uniqueIndex,size:512"`
	Domains []Domain `gorm:"many2many:domain_nameservers;"`
}
