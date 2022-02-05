package api

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/uberswe/domains-sweden/config"
	"gorm.io/gorm"
)

// Controller holds all the variables needed for routes to perform their logic
type Controller struct {
	db     *gorm.DB
	config config.Config
	bundle *i18n.Bundle
}

// New creates a new instance of the routes.Controller
func New(db *gorm.DB, c config.Config, bundle *i18n.Bundle) Controller {
	return Controller{
		db:     db,
		config: c,
		bundle: bundle,
	}
}
