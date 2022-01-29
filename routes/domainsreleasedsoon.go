package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DomainReleaseData struct {
	PageData
	Domains []IndexDomain
	Prev    bool
	Next    bool
	PrevURL string
	NextURL string
}

func (controller Controller) DomainsReleasedSoon(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Domains Being Released Soon")

	page := 1
	perPage := 50

	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	}

	drpd := DomainReleaseData{
		PageData: pd,
	}

	var domains []struct {
		Host       string
		ReleasedAt time.Time
	}

	res := controller.db.Model(&models.Domain{}).
		Select("domains.host, releases.released_at").
		Joins("left join domain_releases on domain_releases.domain_id = domains.id").
		Joins("left join releases on domain_releases.release_id = releases.id").
		Where("releases.released_at > NOW()").
		Order("releases.released_at ASC").
		Offset(perPage * (page - 1)).Limit(perPage).
		Find(&domains)

	if res.Error != nil {
		log.Println(res.Error)
	}

	for _, d := range domains {
		drpd.Domains = append(drpd.Domains, IndexDomain{
			Host:       d.Host,
			URL:        fmt.Sprintf("/domains/%s", d.Host),
			ReleasesAt: d.ReleasedAt.Format("2006-01-02"),
		})
	}

	if len(drpd.Domains) >= perPage {
		drpd.Next = true
		drpd.NextURL = fmt.Sprintf("/domains-released-soon/%d", page+1)
	}
	if page > 1 {
		drpd.Prev = true
		drpd.PrevURL = fmt.Sprintf("/domains-released-soon/%d", page-1)
	}

	c.HTML(http.StatusOK, "domainsreleasedsoon.html", drpd)
}
