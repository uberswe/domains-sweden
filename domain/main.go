package domain

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/uberswe/domains-sweden/models"
	"github.com/uberswe/domains-sweden/sitemap"
	"github.com/uberswe/domains-sweden/text"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	seDomains = "https://data.internetstiftelsen.se/bardate_domains.json"
	nuDomains = "https://data.internetstiftelsen.se/bardate_domains_nu.json"
)

type Service struct {
	DB *gorm.DB
}

type Nameserver struct {
	Domain string
	TTL    int32
}

type Response struct {
	Data    []Domain `json:"data"`
	Domains map[string]Domain
}

type Domain struct {
	Processed   bool
	Name        string `json:"name"`
	ReleaseAt   string `json:"release_at"`
	Nameservers []Nameserver
}

func New(db *gorm.DB) *Service {
	s := Service{
		DB: db,
	}
	return &s
}

func (s *Service) Poll() {
	// Call every 6 hours in a thread
	s.run()
	for range time.Tick(time.Hour * 1) {
		s.run()
	}
}

func (s *Service) run() {
	data := s.load()
	if len(data) > 0 {
		var results []models.Domain
		res := s.DB.Model(&models.Domain{}).Preload("Nameservers").Preload("Releases").FindInBatches(&results, 1000, func(tx *gorm.DB, batch int) error {
			for _, result := range results {
				domainDecoded, _ := idna.ToASCII(result.Host)
				if d, ok := data[domainDecoded]; ok {
					updated := false
					foundRelease := false
					foundNameserver := false
					if d.ReleaseAt != "" {
						parse, err := time.Parse("2006-01-02", d.ReleaseAt)
						if err != nil {
							log.Println(err)
							continue
						}
						for _, release := range result.Releases {
							if release.ReleasedAt.Equal(parse) {
								foundRelease = true
							}
						}
						if !foundRelease {
							release := models.Release{
								ReleasedAt: &parse,
							}
							res := s.DB.Model(&models.Release{}).Where(release).First(&release)
							if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
								log.Println(res.Error)
								return res.Error
							}
							result.Releases = append(result.Releases, release)
							updated = true
						}
					}
					for _, ns := range d.Nameservers {
						found := false
						for _, ns2 := range result.Nameservers {
							if ns.Domain == ns2.Host {
								found = true
								foundNameserver = true
							}
						}
						if !found {
							foundNameserver = false
						}
					}
					if !foundNameserver {
						result.Nameservers = nil
						for _, ns := range d.Nameservers {
							nsDecoded, _ := idna.ToASCII(ns.Domain)
							nameserver := models.Nameserver{
								Host: nsDecoded,
							}
							res := s.DB.Model(&models.Nameserver{}).Where(nameserver).First(&nameserver)
							if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
								log.Println(res.Error)
								return res.Error
							}
							result.Nameservers = append(result.Nameservers, nameserver)
							updated = true
						}
					}
					if updated {
						s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&result)
					}
					d.Processed = true
					data[result.Host] = d
				}
			}
			// returns error will stop future batches
			return nil
		})
		if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
			log.Println(res.Error)
			return
		}

		for domain, obj := range data {
			if !obj.Processed {
				// The domain does not exist in our database
				domainDecoded, _ := idna.ToASCII(domain)
				d := models.Domain{
					Host: domainDecoded,
				}
				if obj.ReleaseAt != "" {
					parse, err := time.Parse("2006-01-02", obj.ReleaseAt)
					if err == nil {
						release := models.Release{
							ReleasedAt: &parse,
						}
						res = s.DB.Where(release).First(&release)
						if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
							log.Println(res.Error)
							return
						}
						d.Releases = append(d.Releases, release)
					} else {
						log.Println(err)
					}
				}
				for _, ns := range obj.Nameservers {
					nsDecoded, _ := idna.ToASCII(ns.Domain)
					nameserver := models.Nameserver{
						Host: nsDecoded,
					}
					res = s.DB.Where(nameserver).First(&nameserver)
					if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
						log.Println(res.Error)
						return
					}
					d.Nameservers = append(d.Nameservers, nameserver)
				}
				s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&d)
			}
		}
	}

	// Aggregate nameservers

	var domainNameservers []struct {
		NameserverID int
		Count        int
	}

	log.Println("aggregating nameserver counts")
	res := s.DB.Table("domain_nameservers").
		Select("domain_nameservers.nameserver_id, COUNT(domain_nameservers.domain_id) AS Count").
		Group("domain_nameservers.nameserver_id").
		Find(&domainNameservers)

	if res.Error != nil {
		log.Println(res.Error)
	}

	for _, dn := range domainNameservers {
		ns := models.NameserverAggregate{
			NameserverID: dn.NameserverID,
		}

		res = s.DB.Where(ns).First(&ns)
		if res.Error != nil {
			log.Println(res.Error)
		}
		ns.Count = dn.Count

		s.DB.Save(&ns)
	}
	// remove duplicate domains
	var domains []struct {
		IDs   string
		Host  string
		Count int64
	}

	log.Println("removing duplicate domains")
	res = s.DB.Raw("SELECT GROUP_CONCAT(id) AS ids, host, COUNT(*) AS count FROM domains GROUP BY host HAVING count > 1").Find(&domains)
	if res.Error != nil {
		log.Println(res.Error)
	}

	for _, d := range domains {
		log.Println("domains", d.IDs, "Host", d.Host, "Count", d.Count)
		idParts := strings.Split(d.IDs, ",")
		for i, p := range idParts {
			if i != 0 || !text.IsASCII(d.Host) {
				s.deleteDomain(p)
			}
		}
	}
	log.Println("generating sitemaps")
	sitemapService := sitemap.New(s.DB)
	sitemapService.GenerateAllSitemaps()
	log.Println("hourly update complete")
}

func (s *Service) deleteDomain(domainID string) {
	dModel := models.Domain{}
	if pi, err := strconv.Atoi(domainID); err == nil {
		dModel.ID = uint(pi)
		err = s.DB.Model(&dModel).Association("Nameservers").Clear()
		if err != nil {
			log.Println(err)
		}
		err = s.DB.Model(&dModel).Association("Releases").Clear()
		if err != nil {
			log.Println(err)
		}
		err = s.DB.Model(&dModel).Association("Parses").Clear()
		if err != nil {
			log.Println(err)
		}
		res := s.DB.Unscoped().Delete(&dModel)
		if res.Error != nil {
			log.Println(res.Error)
		}
	} else {
		log.Println("domain duplication filtering failed", err)
	}
}

func (s *Service) load() map[string]Domain {
	var domains map[string]Domain
	var fetch models.Fetch
	res := s.DB.Order("created_at DESC").First(&fetch)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Println(res.Error)
		return nil
	}

	if fetch.ID == 0 || time.Now().Add(-6*time.Hour).After(fetch.CreatedAt) {
		var fetchNew models.Fetch
		data := loadExpiringDomains("se")
		fetchNew.ReleasingSEDomains = len(data.Data)
		nuData := loadExpiringDomains("nu")
		fetchNew.ReleasingNUDomains = len(data.Data)
		data.Data = append(data.Data, nuData.Data...)
		data.Domains = loadActiveDomains("se")
		fetchNew.ActiveSEDomains = len(data.Domains)
		nuDomainMap := loadActiveDomains("nu")
		fetchNew.ActiveNUDomains = len(nuDomainMap)
		for i, d := range nuDomainMap {
			data.Domains[i] = d
		}
		log.Println("Loaded domains from iis.se", len(data.Domains), len(data.Data))
		fetchNew.ActiveDomains = fetchNew.ActiveSEDomains + fetchNew.ActiveNUDomains
		fetchNew.ReleasingDomains = fetchNew.ReleasingSEDomains + fetchNew.ReleasingNUDomains
		res = s.DB.Save(&fetchNew)
		if res.Error != nil {
			log.Println(res.Error)
			return nil
		}

		domains = data.Domains

		if domains == nil {
			domains = make(map[string]Domain)
		}

		for _, d := range data.Data {
			host, _ := idna.ToASCII(strings.TrimRight(d.Name, "."))
			if tmpD, ok := domains[host]; ok {
				tmpD.ReleaseAt = d.ReleaseAt
				domains[host] = tmpD
			} else {
				domains[host] = d
			}
		}
	}

	return domains
}

func loadActiveDomains(segment string) map[string]Domain {
	// dig @zonedata.iis.se se AXFR > se.zone.txt
	// dig @zonedata.iis.se nu AXFR > nu.zone.txt
	domains := make(map[string]Domain)
	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(fmt.Sprintf("%s.", segment))
	ch, err := t.In(m, "zonedata.iis.se:53")
	if err != nil {
		log.Println(err)
		return domains
	}
	if err != nil {
		log.Println(err)
		return domains
	}
	for env := range ch {
		if env.Error != nil {
			err = env.Error
			break
		}
		for _, rr := range env.RR {
			switch v := rr.(type) {
			case *dns.NS:
				dn, _ := idna.ToUnicode(strings.TrimRight(v.Hdr.Name, "."))
				if _, ok := domains[dn]; !ok {
					domains[dn] = Domain{}
				}
				dTmp := domains[dn]
				ns, _ := idna.ToUnicode(strings.TrimRight(v.Ns, "."))
				dTmp.Nameservers = append(dTmp.Nameservers, Nameserver{
					Domain: ns,
				})
				domains[dn] = dTmp
			}
		}
	}
	if err != nil {
		log.Println(err)
		return domains
	}
	return domains
}

func loadExpiringDomains(segment string) (data Response) {
	var err error
	var req *http.Request
	log.Println("loading from url")
	client := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	if segment == "se" {
		req, err = http.NewRequest(http.MethodGet, seDomains, nil)
	} else if segment == "nu" {
		req, err = http.NewRequest(http.MethodGet, nuDomains, nil)
	}
	if err != nil {
		log.Println(err)
		return data
	}

	req.Header.Set("User-Agent", "domäner.xyz parser, contact web@domaner.xyz in case of abuse. Also available on Github https://github.com/uberswe/domains-sweden")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Println(getErr)
		return data
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(readErr)
		return data
	}

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Println(jsonErr)
		return data
	}
	return data
}

func Title(s string) string {
	hash, _ := idna.ToUnicode(s)
	parts := strings.Split(hash, ".")
	parts[0] = strings.Title(strings.ToLower(parts[0]))
	return strings.Join(parts, ".")
}

func ToUnicode(s string) string {
	hash, _ := idna.ToUnicode(s)
	return hash
}
