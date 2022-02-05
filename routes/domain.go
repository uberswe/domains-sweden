package routes

import (
	"github.com/gin-gonic/gin"
	domainservice "github.com/uberswe/domains-sweden/domain"
	"github.com/uberswe/domains-sweden/models"
	"golang.org/x/net/idna"
	"net/http"
	"time"
)

type DomainPageData struct {
	PageData
	FirstSeen      string
	HasNameservers bool
	Nameservers    []models.Nameserver
	HasReleaseAt   bool
	ReleaseAt      string
}

func (controller Controller) Domain(c *gin.Context) {
	dpd := controller.DefaultPageData(c)

	domain := c.Param("domain")

	domainModel := models.Domain{
		Host: domain,
	}

	res := controller.db.Where(domainModel).Preload("Nameservers").Preload("Releases").First(&domainModel)
	if res.Error != nil {
		c.HTML(http.StatusNotFound, "404.html", dpd)
		return
	}

	hash, _ := idna.ToUnicode(domainModel.Host)

	dpd.Title = domainservice.Title(hash)
	pd := DomainPageData{
		PageData:    dpd,
		FirstSeen:   domainModel.CreatedAt.Format("2006-01-02"),
		Nameservers: domainModel.Nameservers,
	}
	if len(pd.Nameservers) > 0 {
		pd.HasNameservers = true
	}
	for _, r := range domainModel.Releases {
		if r.ReleasedAt != nil && r.ReleasedAt.After(time.Now()) {
			pd.HasReleaseAt = true
			pd.ReleaseAt = r.ReleasedAt.Format("2006-01-02")
			break
		}
	}

	c.HTML(http.StatusOK, "domain.html", pd)
}
