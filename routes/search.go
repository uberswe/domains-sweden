package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/domain"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SearchData holds additional data needed to render the search HTML page
type SearchData struct {
	PageData
	Results        []SearchDomain
	Prev           bool
	Next           bool
	PrevURL        string
	NextURL        string
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

type SearchDomain struct {
	Host       string
	URL        string
	ReleasedAt string
	Requested  string
	Status     string
}

type SearchQueryResult struct {
	Host       string
	ReleasedAt *time.Time
	Requested  *time.Time
	LoadTime   *float64
}

// Search renders the search HTML page and any search results
func (controller Controller) Search(c *gin.Context) {
	page := 1
	resultsPerPage := 50
	pdS := controller.DefaultPageData(c)
	pdS.Title = pdS.Trans("Search")
	pd := SearchData{
		PageData: pdS,
	}
	search := ""
	if c.Request.Method == "POST" && c.Request.RequestURI == "/search" {
		search = c.PostForm("search")
		min := c.PostForm("min")
		max := c.PostForm("max")
		extension := c.PostForm("extension")
		website := c.PostForm("website")
		releasing := c.PostForm("releasing")
		expiring := c.PostForm("expiring")
		nospecialchars := c.PostForm("nospecialchars")
		nonumbers := c.PostForm("nonumbers")
		c.Redirect(http.StatusTemporaryRedirect, buildSearchURL(search, min, max, extension, website, releasing, expiring, nospecialchars, nonumbers, 1))
		return
	} else {
		search = c.Param("query")
		if i, err := strconv.Atoi(c.Param("page")); err == nil {
			page = i
		}
	}
	min := c.Query("min")
	max := c.Query("max")
	extension := c.Query("extension")
	website := c.Query("website")
	releasing := c.Query("releasing")
	expiring := c.Query("expiring")
	nospecialchars := c.Query("nospecialchars")
	nonumbers := c.Query("nonumbers")

	search, _ = idna.ToASCII(search)

	var results []SearchQueryResult

	searchFilter := fmt.Sprintf("%s%s%s", "%", search, "%")
	search2 := fmt.Sprintf("%s%s", "%", search)
	search4 := fmt.Sprintf("%s%s", search, "%")

	q := controller.db.
		Table("domains").
		Select("domains.host, parses.requested, parses.load_time, releases.released_at").
		Joins("LEFT JOIN parses ON parses.domain_id = domains.id").
		Joins("LEFT JOIN domain_releases ON domain_releases.domain_id = domains.id").
		Joins("LEFT JOIN releases ON domain_releases.release_id = releases.id").
		Order("LENGTH(domains.host)").
		Limit(resultsPerPage).
		Offset(resultsPerPage * (page - 1))

	if search != "" {
		q = q.Where("BINARY domains.host LIKE ?", strings.ToLower(searchFilter)).
			Order(gorm.Expr("CASE WHEN BINARY domains.host LIKE ? THEN 1 WHEN BINARY domains.host LIKE ? THEN 2 WHEN BINARY domains.host LIKE ? THEN 4 ELSE 3 END", strings.ToLower(search), strings.ToLower(search2), strings.ToLower(search4)))
	}

	if i, err := strconv.Atoi(min); err == nil {
		// Add +3 because .se or .nu is not counted in the length
		q = q.Where("LENGTH(domains.host) >= ?", i+3)
	}

	if i, err := strconv.Atoi(max); err == nil {
		q = q.Where("LENGTH(domains.host) <= ?", i+3)
	}

	if extension != "" {
		q = q.Where("domains.host LIKE ?", fmt.Sprintf("%s%s", "%", extension))
	}

	if website == "1" {
		q = q.Where("parses.load_time > ?", 0)
	}

	if website == "2" {
		q = q.Where(controller.db.Where("parses.id IS NULL").
			Or("parses.load_time = ?", 0))
	}

	if i, err := strconv.Atoi(releasing); err == nil {
		from := time.Now().Add(time.Hour * time.Duration(i) * 24)
		q = q.Where("releases.released_at < ?", from).Where("releases.released_at > NOW()")
	}

	if expiring != "" {
		q = q.Where("releases.released_at IS NOT NULL").Where("releases.released_at > NOW()")
	}

	if nospecialchars != "" {
		q = q.Where("domains.host NOT REGEXP ?", "[-]")
	}

	if nonumbers != "" {
		q = q.Where("domains.host NOT REGEXP ?", "[0-9]")
	}

	res := q.Find(&results)
	if res.Error != nil || len(results) == 0 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: pdS.Trans("No results found"),
		})
		log.Println(res.Error)
		c.HTML(http.StatusOK, "search.html", pd)
		return
	}

	for i := range results {
		host, _ := idna.ToUnicode(results[i].Host)
		status := pd.Trans("Unknown")
		if results[i].LoadTime != nil {
			if *results[i].LoadTime > 0 {
				status = pd.Trans("OK")
			} else {
				status = pd.Trans("Error")
			}
		}
		releasedAt := ""
		if results[i].ReleasedAt != nil {
			releasedAt = results[i].ReleasedAt.Format("2006-01-02")
		}
		requested := ""
		if results[i].Requested != nil {
			requested = results[i].Requested.Format("2006-01-02")
		}
		pd.Results = append(pd.Results, SearchDomain{
			Host:       domain.Title(host),
			URL:        fmt.Sprintf("/domains/%s", domain.ToUnicode(results[i].Host)),
			ReleasedAt: releasedAt,
			Requested:  requested,
			Status:     status,
		})
	}

	if len(pd.Results) >= resultsPerPage {
		pd.Next = true
		pd.NextURL = buildSearchURL(search, min, max, extension, website, releasing, expiring, nospecialchars, nonumbers, page+1)
	}
	if page > 1 {
		pd.Prev = true
		pd.PrevURL = buildSearchURL(search, min, max, extension, website, releasing, expiring, nospecialchars, nonumbers, page-1)
	}

	pd.Count = make(map[int]string)
	for i := 1; i <= 90; i++ {
		pd.Count[i] = fmt.Sprintf("%d", i)
	}

	pd.Search = search
	pd.Min = min
	pd.Max = max
	pd.Extension = extension
	pd.Website = website
	pd.Releasing = releasing
	pd.Expiring = expiring
	pd.NoSpecialChars = nospecialchars
	pd.NoNumbers = nonumbers

	c.HTML(http.StatusOK, "search.html", pd)
}

func buildSearchURL(search, min, max, extension, website, releasing, expiring, nospecialchars, nonumbers string, page int) string {
	params := url.Values{}
	if min != "" {
		params.Add("min", min)
	}
	if max != "" {
		params.Add("max", max)
	}
	if extension != "" {
		params.Add("extension", extension)
	}
	if website != "" {
		params.Add("website", website)
	}
	if releasing != "" {
		params.Add("releasing", releasing)
	}
	if expiring != "" {
		params.Add("expiring", expiring)
	}
	if nospecialchars != "" {
		params.Add("nospecialchars", nospecialchars)
	}
	if nonumbers != "" {
		params.Add("nonumbers", nonumbers)
	}

	encodedParams := params.Encode()
	if encodedParams != "" {
		search = url.QueryEscape(search) + "?" + encodedParams
	}
	return fmt.Sprintf("/search/%d/%s", page, search)
}
