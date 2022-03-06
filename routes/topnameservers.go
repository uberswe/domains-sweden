package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	domainservice "github.com/uberswe/domains-sweden/domain"
	"github.com/uberswe/domains-sweden/models"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type TopNameserversData struct {
	PageData
	Nameservers []IndexNameserver
	Prev        bool
	Next        bool
	PrevURL     string
	NextURL     string
}

type TopNameserverCache struct {
	Cached time.Time
	Data   TopNameserversData
}

var topNameserverDataCache = map[int]TopNameserverCache{}

func (controller Controller) TopNameservers(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Top Nameservers")

	page := 1
	perPage := 50

	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	}

	tnpd := TopNameserversData{}

	var domainNameservers []models.NameserverAggregate

	if _, ok := topNameserverDataCache[page]; page > 10 || !ok || topNameserverDataCache[page].Cached.Before(time.Now().Add(-6*time.Hour)) {

		res := controller.db.Model(models.NameserverAggregate{}).
			Order("count DESC").
			Offset(perPage * (page - 1)).Limit(perPage).
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
			tnpd.Nameservers = append(tnpd.Nameservers, IndexNameserver{
				Host:  domainservice.ToUnicode(ns.Host),
				URL:   fmt.Sprintf("/nameservers/%s", domainservice.ToUnicode(ns.Host)),
				Count: int64(count),
			})
		}

		sort.Slice(tnpd.Nameservers, func(i, j int) bool {
			return tnpd.Nameservers[i].Count > tnpd.Nameservers[j].Count
		})

		if len(tnpd.Nameservers) >= perPage {
			tnpd.Next = true
			tnpd.NextURL = fmt.Sprintf("/top-nameservers/%d", page+1)
		}
		if page > 1 {
			tnpd.Prev = true
			tnpd.PrevURL = fmt.Sprintf("/top-nameservers/%d", page-1)
		}

		topNameserverDataCache[page] = TopNameserverCache{
			Cached: time.Now(),
			Data:   tnpd,
		}
	} else {
		tnpd = topNameserverDataCache[page].Data
	}

	tnpd.PageData = pd

	c.HTML(http.StatusOK, "topnameservers.html", tnpd)
}
