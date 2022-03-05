package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	domainservice "github.com/uberswe/domains-sweden/domain"
	"github.com/uberswe/domains-sweden/models"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"html/template"
	"log"
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
	HasParse       bool
	Screenshot     template.HTML
	PageSize       string
	LoadTime       string
	ParsedAt       string
}

func (controller Controller) Domain(c *gin.Context) {
	dpd := controller.DefaultPageData(c)

	domain := c.Param("domain")

	domainDecoded, _ := idna.ToASCII(domain)

	domainModel := models.Domain{
		Host: domainDecoded,
	}

	res := controller.db.Where(domainModel).Preload("Nameservers").Preload("Releases").First(&domainModel)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			domainModel = models.Domain{
				Host: domain,
			}
			res = controller.db.Where(domainModel).Preload("Nameservers").Preload("Releases").First(&domainModel)
			if res.Error != nil {
				if res.Error == gorm.ErrRecordNotFound {
					// TODO perform a dns lookup to add the domain to our database
				}
				c.HTML(http.StatusNotFound, "404.html", dpd)
				return
			}
		} else {
			c.HTML(http.StatusNotFound, "404.html", dpd)
			return
		}
	}

	go controller.parser.Parse(domainModel)

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

	parse := models.Parse{
		DomainID: domainModel.ID,
	}

	res = controller.db.Where(parse).Order("created_at DESC").First(&parse)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Println(res.Error)
	}

	if parse.ID > 0 {
		pd.HasParse = true
		if parse.Error != nil {
			pd.LoadTime = pd.Trans("Error loading page")
		} else {
			pd.Screenshot = template.HTML(fmt.Sprintf("<img src=\"/domains/%s/screenshot.jpg\" class=\"img-thumbnail mx-auto d-block\" alt=\"%s\">", domainModel.Host, pd.Title))
			pd.LoadTime = fmt.Sprintf("%0.3f %s", parse.LoadTime, pd.Trans("Seconds"))
			pd.PageSize = fmt.Sprintf("%0.2f %s", parse.Size, pd.Trans("Mb"))
		}
		pd.ParsedAt = parse.Requested.Format("2006-01-02")
	}

	c.HTML(http.StatusOK, "domain.html", pd)
}

func (controller Controller) DomainScreenshot(c *gin.Context) {
	// Return JPG from remote server
	domain := c.Param("domain")

	domainDecoded, _ := idna.ToASCII(domain)

	domainModel := models.Domain{
		Host: domainDecoded,
	}

	res := controller.db.Where(domainModel).First(&domainModel)
	if res.Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	parseModel := models.Parse{
		DomainID: domainModel.ID,
	}

	res = controller.db.Where(parseModel).Order("created_at DESC").First(&parseModel)
	if res.Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if parseModel.BlurredScreenshotHash != "" {
		fetch, err := controller.sftpService.Fetch(fmt.Sprintf("/screenshots/blurred-%s.jpg", parseModel.BlurredScreenshotHash))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.Data(http.StatusOK, "image/jpeg", fetch)
		return
	}

	c.Status(http.StatusNotFound)
	return
}
