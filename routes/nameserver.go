package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/domain"
	"github.com/uberswe/domains-sweden/models"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type NameserverPageData struct {
	PageData
	FirstSeen  string
	HasDomains bool
	Domains    []SearchDomain
	Prev       bool
	Next       bool
	PrevURL    string
	NextURL    string
}

func (controller Controller) Nameserver(c *gin.Context) {
	dpd := controller.DefaultPageData(c)

	nameserver := c.Param("nameserver")

	nameserverModel := models.Nameserver{
		Host: nameserver,
	}
	perPage := 20
	page := 1
	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	}

	res := controller.db.Where(nameserverModel).First(&nameserverModel)
	if res.Error != nil {
		c.HTML(http.StatusNotFound, "404.html", dpd)
		return
	}

	hash, _ := idna.ToUnicode(nameserverModel.Host)

	dpd.Title = domain.Title(hash)
	pd := NameserverPageData{
		PageData:  dpd,
		FirstSeen: nameserverModel.CreatedAt.Format("2006-01-02"),
	}

	var domains []models.Domain

	err := controller.db.Model(&nameserverModel).Order("id ASC").Offset(perPage * (page - 1)).Limit(perPage).Association("Domains").Find(&domains)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.HTML(http.StatusNotFound, "404.html", dpd)
		return
	}

	for _, d := range domains {
		pd.Domains = append(pd.Domains, SearchDomain{
			Host: domain.Title(d.Host),
			URL:  fmt.Sprintf("/domains/%s", d.Host),
		})
	}

	pd.HasDomains = len(pd.Domains) > 0

	if len(pd.Domains) >= perPage {
		pd.Next = true
		pd.NextURL = fmt.Sprintf("/nameservers/%s/%d", nameserverModel.Host, page+1)
	}
	if page > 1 {
		pd.Prev = true
		pd.PrevURL = fmt.Sprintf("/nameservers/%s/%d", nameserverModel.Host, page-1)
	}

	c.HTML(http.StatusOK, "nameserver.html", pd)
}
