package whois

import (
	"github.com/likexian/whois"
	"gorm.io/gorm"
	"log"
	"time"
)

type Service struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Service {
	s := Service{
		DB: db,
	}
	return &s
}

type Whois struct {
	Raw     string    `json:"raw"`
	Fetched time.Time `json:"fetched"`
}

// Domain performs a whois lookup
func (s *Service) Domain(host string) (*Whois, error) {
	// We don't store any whois data as internetstiftelsen.se does not allow us to store or use any whois query data for .se and .nu domains
	result, err := whois.Whois(host)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Whois{
		Raw:     result,
		Fetched: time.Now(),
	}, nil
}
