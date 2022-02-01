package models

import (
	"gorm.io/gorm"
	"time"
)

type Nameserver struct {
	gorm.Model
	Host                 string   `gorm:"uniqueIndex,size:512"`
	Domains              []Domain `gorm:"many2many:domain_nameservers;"`
	NameserverAggregates []NameserverAggregate
}

type NameserverAggregate struct {
	Count        int
	NameserverID int `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
