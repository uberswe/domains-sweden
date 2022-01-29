package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	var nameservers []struct {
		Host  string
		Count int
	}

	if _, ok := topNameserverDataCache[page]; page > 10 || !ok || topNameserverDataCache[page].Cached.Before(time.Now().Add(-6*time.Hour)) {

		res := controller.db.Table("domain_nameservers").
			Select("nameservers.host, COUNT(domain_nameservers.domain_id) AS count").
			Joins("left join nameservers on domain_nameservers.nameserver_id = nameservers.id").
			Order("COUNT(domain_nameservers.domain_id) DESC").
			Offset(perPage * (page - 1)).Limit(perPage).
			Group("nameservers.host").
			Find(&nameservers)

		if res.Error != nil {
			log.Println(res.Error)
		}

		for _, ns := range nameservers {
			tnpd.Nameservers = append(tnpd.Nameservers, IndexNameserver{
				Host:  ns.Host,
				URL:   fmt.Sprintf("/nameservers/%s", ns.Host),
				Count: fmt.Sprintf("%d", ns.Count),
			})
		}

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
