package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/models"
	"log"
	"net/http"
	"strconv"
)

func (controller Controller) Sitemap(c *gin.Context) {
	var domainCount int64
	controller.db.Model(&models.Domain{}).Count(&domainCount)
	var nameserverCount int64
	controller.db.Model(&models.Nameserver{}).Count(&nameserverCount)

	sitemap := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	sitemap += "<sitemapindex xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"
	sitemap += "	<sitemap>"
	sitemap += "		<loc>https://www.xn--domner-dua.xyz/sitemap/default/sitemap.xml</loc>"
	sitemap += "	</sitemap>"
	for i := int64(0); i*1000 < domainCount; i++ {
		sitemap += "	<sitemap>"
		sitemap += fmt.Sprintf("		<loc>https://www.xn--domner-dua.xyz/sitemap/domains/%d/sitemap.xml</loc>", i+1)
		sitemap += "	</sitemap>"
	}
	for i := int64(0); i*1000 < nameserverCount; i++ {
		sitemap += "	<sitemap>"
		sitemap += fmt.Sprintf("		<loc>https://www.xn--domner-dua.xyz/sitemap/nameservers/%d/sitemap.xml</loc>", i+1)
		sitemap += "	</sitemap>"
	}
	sitemap += "</sitemapindex>"

	c.Data(http.StatusOK, "application/xml", []byte(sitemap))
}

func (controller Controller) SitemapDefault(c *gin.Context) {
	sitemap := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	sitemap += "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"
	sitemap += "	<url>"
	sitemap += "		<loc>https://www.xn--domner-dua.xyz/</loc>"
	sitemap += "		<changefreq>daily</changefreq>"
	sitemap += "		<priority>0.9</priority>"
	sitemap += "	</url>"
	sitemap += "	<url>"
	sitemap += "		<loc>https://www.xn--domner-dua.xyz/domains-released-soon</loc>"
	sitemap += "		<changefreq>daily</changefreq>"
	sitemap += "		<priority>0.8</priority>"
	sitemap += "	</url>"
	sitemap += "	<url>"
	sitemap += "		<loc>https://www.xn--domner-dua.xyz/top-nameservers</loc>"
	sitemap += "		<changefreq>daily</changefreq>"
	sitemap += "		<priority>0.8</priority>"
	sitemap += "	</url>"
	sitemap += "</urlset>"
	c.Data(http.StatusOK, "application/xml", []byte(sitemap))
}

func (controller Controller) SitemapDomains(c *gin.Context) {

	page := 1
	perPage := 1000

	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	} else {
		c.Status(http.StatusNotFound)
		return
	}

	var domains []models.Domain

	res := controller.db.Model(&models.Domain{}).Select("host, updated_at").Order("id ASC").Offset(perPage * (page - 1)).Limit(perPage).Find(&domains)
	if res.Error != nil {
		log.Println(res.Error)
		c.Status(http.StatusNotFound)
		return
	}

	sitemap := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	sitemap += "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"
	for _, d := range domains {
		sitemap += "	<url>"
		sitemap += fmt.Sprintf("		<loc>https://www.xn--domner-dua.xyz/domains/%s</loc>", d.Host)
		sitemap += "		<changefreq>monthly</changefreq>"
		sitemap += fmt.Sprintf("        <lastmod>%s</lastmod>", d.UpdatedAt.Format("2006-01-02"))
		sitemap += "		<priority>0.7</priority>"
		sitemap += "	</url>"
	}
	sitemap += "</urlset>"
	c.Data(http.StatusOK, "application/xml", []byte(sitemap))
}

func (controller Controller) SitemapNameservers(c *gin.Context) {
	page := 1
	perPage := 1000

	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	} else {
		c.Status(http.StatusNotFound)
		return
	}

	var nameservers []models.Nameserver

	res := controller.db.Model(&models.Nameserver{}).Select("host, updated_at").Order("id ASC").Offset(perPage * (page - 1)).Limit(perPage).Find(&nameservers)
	if res.Error != nil {
		log.Println(res.Error)
		c.Status(http.StatusNotFound)
		return
	}

	sitemap := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	sitemap += "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"
	for _, ns := range nameservers {
		sitemap += "	<url>"
		sitemap += fmt.Sprintf("		<loc>https://www.xn--domner-dua.xyz/nameservers/%s</loc>", ns.Host)
		sitemap += "		<changefreq>monthly</changefreq>"
		sitemap += fmt.Sprintf("        <lastmod>%s</lastmod>", ns.UpdatedAt.Format("2006-01-02"))
		sitemap += "		<priority>0.6</priority>"
		sitemap += "	</url>"
	}
	sitemap += "</urlset>"
	c.Data(http.StatusOK, "application/xml", []byte(sitemap))
}
