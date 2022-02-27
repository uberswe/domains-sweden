package models

import (
	"github.com/chromedp/cdproto/debugger"
	"gorm.io/gorm"
	"time"
)

type Parse struct {
	gorm.Model
	debugger.ScriptLanguage
	ContentHash           string  `gorm:"type:text"`
	ScreenshotHash        string  `gorm:"type:text"`
	BlurredScreenshotHash string  `gorm:"type:text"`
	Content               []byte  `gorm:"type:longblob"`
	Screenshot            []byte  `gorm:"type:blob"`
	BlurredScreenshot     []byte  `gorm:"type:blob"` // There is a lot of porn and other weird stuff out there, so I have opted to blur the screenshots. It is enough to give an idea of what the site was in the past.
	Size                  float64 // in Mb
	LoadTime              float64 // in Seconds
	Error                 *string
	Requested             time.Time
	Events                []ParseEvent
	DomainID              uint
	Domain                Domain
}
