package sitemap

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/models"
	"gorm.io/gorm"
	"log"
	"strconv"
)

var perPage = 1000

type Service struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Service {
	s := Service{
		DB: db,
	}
	return &s
}

func (s *Service) GenerateAllSitemaps() {
	s.Main(true)
	page := 1
	identifier := fmt.Sprintf("domains_%d", page)
	for {
		content, err := s.sitemapFromDomains(page, identifier)
		if err != nil {
			break
		}
		if content == "" {
			break
		}
		page++
		identifier = fmt.Sprintf("domains_%d", page)
	}
	page = 1
	for {
		content, err := s.sitemapFromNameservers(page, identifier)
		if err != nil {
			break
		}
		if content == "" {
			break
		}
		page++
		identifier = fmt.Sprintf("nameservers_%d", page)
	}
}

func (s *Service) Default() string {
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
	return sitemap
}

func (s *Service) Domains(c *gin.Context, regenerate bool) (string, error) {
	page := 1
	identifier := fmt.Sprintf("domains_%d", page)

	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	} else {
		return "", errors.New("not found")
	}

	if !regenerate {
		sitemap, exists := s.fetchSitemap(identifier)
		if exists {
			return sitemap, nil
		}
	}

	return s.sitemapFromDomains(page, identifier)
}

func (s *Service) sitemapFromDomains(page int, identifier string) (string, error) {
	var domains []models.Domain

	res := s.DB.Model(&models.Domain{}).Select("host, updated_at").Order("id ASC").Offset(perPage * (page - 1)).Limit(perPage).Find(&domains)
	if res.Error != nil {
		return "", res.Error
	}
	if len(domains) == 0 {
		return "", errors.New("no domains found")
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
	go func() {
		err := s.storeSitemap(identifier, sitemap)
		if err != nil {
			log.Println(err)
		}
	}()
	return sitemap, nil
}

func (s *Service) Nameservers(c *gin.Context, regenerate bool) (string, error) {
	page := 1
	identifier := fmt.Sprintf("nameservers_%d", page)

	if i, err := strconv.Atoi(c.Param("page")); err == nil {
		page = i
	} else {
		return "", errors.New("not found")
	}

	if !regenerate {
		sitemap, exists := s.fetchSitemap(identifier)
		if exists {
			return sitemap, nil
		}
	}

	return s.sitemapFromNameservers(page, identifier)
}

func (s *Service) sitemapFromNameservers(page int, identifier string) (string, error) {
	var nameservers []models.Nameserver

	res := s.DB.Model(&models.Nameserver{}).Select("host, updated_at").Order("id ASC").Offset(perPage * (page - 1)).Limit(perPage).Find(&nameservers)
	if res.Error != nil {
		return "", res.Error
	}

	if len(nameservers) == 0 {
		return "", errors.New("no domains found")
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
	go func() {
		err := s.storeSitemap(identifier, sitemap)
		if err != nil {
			log.Println(err)
		}
	}()
	return sitemap, nil
}

func (s *Service) Main(regenerate bool) (sitemap string) {
	var exists bool
	if !regenerate {
		sitemap, exists = s.fetchSitemap("main")
		if exists {
			return sitemap
		}
	}

	var domainCount int64
	s.DB.Model(&models.Domain{}).Count(&domainCount)
	var nameserverCount int64
	s.DB.Model(&models.Nameserver{}).Count(&nameserverCount)

	sitemap = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
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
	go func() {
		err := s.storeSitemap("main", sitemap)
		if err != nil {
			log.Println(err)
		}
	}()
	return sitemap
}

func (s *Service) fetchSitemap(identifier string) (content string, exists bool) {
	sm := models.Sitemap{
		Identifier: identifier,
	}
	res := s.DB.Where(sm).First(&sm)
	if res.Error != nil {
		log.Println(res.Error)
		return "", false
	}
	return sm.Content, true
}

func (s *Service) storeSitemap(identifier string, content string) error {
	if _, exists := s.fetchSitemap(identifier); !exists {
		sm := models.Sitemap{
			Identifier: identifier,
			Content:    content,
		}
		res := s.DB.Save(&sm)
		if res.Error != nil {
			log.Println(res.Error)
			return res.Error
		}
		return nil
	}

	sm := models.Sitemap{
		Identifier: identifier,
	}

	res := s.DB.Where(sm).First(&sm)
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}

	sm.Content = content

	res = s.DB.Save(&sm)
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}
	return nil
}
