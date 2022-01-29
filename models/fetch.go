package models

import (
	"gorm.io/gorm"
)

type Fetch struct {
	gorm.Model
	ActiveDomains      int
	ReleasingDomains   int
	ActiveSEDomains    int
	ActiveNUDomains    int
	ReleasingSEDomains int
	ReleasingNUDomains int
}
