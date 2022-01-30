package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/models"
	"log"
	"net/http"
	"sort"
	"time"
)

type IndexData struct {
	PageData
	Domains     []IndexDomain
	Nameservers []IndexNameserver
}

type IndexDomain struct {
	Host       string
	URL        string
	ReleasesAt string
}

type IndexNameserver struct {
	Host  string
	URL   string
	Count int64
}

type IndexCache struct {
	Cached    time.Time
	IndexData IndexData
}

var indexCache = IndexCache{}

// Index renders the HTML of the index page
func (controller Controller) Index(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Home")

	ipd := IndexData{}

	if indexCache.Cached.Before(time.Now().Add(-6 * time.Hour)) {

		var domains []struct {
			Host       string
			ReleasedAt time.Time
		}

		res := controller.db.Model(&models.Domain{}).
			Select("domains.host, releases.released_at").
			Joins("left join domain_releases on domain_releases.domain_id = domains.id").
			Joins("left join releases on domain_releases.release_id = releases.id").
			Where("releases.released_at > NOW()").
			Order("domains.id ASC").
			Order("releases.released_at ASC").
			Offset(50).
			Limit(20).
			Find(&domains)

		if res.Error != nil {
			log.Println(res.Error)
		}

		for _, d := range domains {
			ipd.Domains = append(ipd.Domains, IndexDomain{
				Host:       d.Host,
				URL:        fmt.Sprintf("/domains/%s", d.Host),
				ReleasesAt: d.ReleasedAt.Format("2006-01-02"),
			})
		}

		var domainNameservers []struct {
			NameserverID int
			Count        int
		}

		res = controller.db.Table("domain_nameservers").
			Select("domain_nameservers.nameserver_id, COUNT(domain_nameservers.domain_id) AS count").
			Order("COUNT(domain_nameservers.domain_id) DESC").
			Limit(20).
			Group("domain_nameservers.nameserver_id").
			Find(&domainNameservers)

		if res.Error != nil {
			log.Println(res.Error)
		}

		var nameservers []models.Nameserver

		var nsIds []int

		for _, dns := range domainNameservers {
			nsIds = append(nsIds, dns.NameserverID)
		}

		controller.db.Find(&nameservers, nsIds)

		for _, ns := range nameservers {
			count := 0
			for _, dn := range domainNameservers {
				if dn.NameserverID == int(ns.ID) {
					count = dn.Count
				}
			}
			ipd.Nameservers = append(ipd.Nameservers, IndexNameserver{
				Host:  ns.Host,
				URL:   fmt.Sprintf("/nameservers/%s", ns.Host),
				Count: int64(count),
			})
		}
		sort.Slice(ipd.Nameservers, func(i, j int) bool {
			return ipd.Nameservers[i].Count > ipd.Nameservers[j].Count
		})
		indexCache.IndexData = ipd
		indexCache.Cached = time.Now()
	} else {
		ipd = indexCache.IndexData
	}
	ipd.PageData = pd

	c.HTML(http.StatusOK, "index.html", ipd)
}
