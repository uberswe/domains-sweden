package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/sitemap"
	"net/http"
)

func (controller Controller) Sitemap(c *gin.Context) {
	sitemapService := sitemap.New(controller.db)
	s := sitemapService.Main(false)
	c.Data(http.StatusOK, "application/xml", []byte(s))
}

func (controller Controller) SitemapDefault(c *gin.Context) {
	sitemapService := sitemap.New(controller.db)
	s := sitemapService.Default()
	c.Data(http.StatusOK, "application/xml", []byte(s))
}

func (controller Controller) SitemapDomains(c *gin.Context) {
	sitemapService := sitemap.New(controller.db)
	s, err := sitemapService.Domains(c, false)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, "application/xml", []byte(s))
}

func (controller Controller) SitemapNameservers(c *gin.Context) {
	sitemapService := sitemap.New(controller.db)
	s, err := sitemapService.Nameservers(c, false)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, "application/xml", []byte(s))
}
