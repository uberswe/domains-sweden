package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/domain"
	"github.com/uberswe/domains-sweden/models"
	"log"
	"net/http"
	"sort"
	"time"
)

type IndexData struct {
	PageData
	Domains        []IndexDomain
	Nameservers    []IndexNameserver
	Search         string
	Min            string
	Max            string
	Extension      string
	Website        string
	Releasing      string
	Expiring       string
	NoSpecialChars string
	NoNumbers      string
	Count          map[int]string
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

		res := controller.db.Model(&models.Release{}).
			Select("domains.host, releases.released_at").
			Joins("left join domain_releases on domain_releases.release_id = releases.id").
			Joins("left join domains on domain_releases.domain_id = domains.id").
			Where("releases.released_at > NOW()").
			Order("domains.id ASC").
			Offset(50).
			Limit(20).
			Find(&domains)

		if res.Error != nil {
			log.Println(res.Error)
		}

		for _, d := range domains {
			ipd.Domains = append(ipd.Domains, IndexDomain{
				Host:       domain.Title(d.Host),
				URL:        fmt.Sprintf("/domains/%s", domain.ToUnicode(d.Host)),
				ReleasesAt: d.ReleasedAt.Format("2006-01-02"),
			})
		}

		var domainNameservers []models.NameserverAggregate

		res = controller.db.Model(models.NameserverAggregate{}).
			Order("count DESC").
			Limit(20).
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
				Host:  domain.ToUnicode(ns.Host),
				URL:   fmt.Sprintf("/nameservers/%s", domain.ToUnicode(ns.Host)),
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

	ipd.Count = make(map[int]string)
	for i := 1; i <= 90; i++ {
		ipd.Count[i] = fmt.Sprintf("%d", i)
	}

	c.HTML(http.StatusOK, "index.html", ipd)
}
