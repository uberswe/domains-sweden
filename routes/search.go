package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/models"
	"golang.org/x/net/idna"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// SearchData holds additional data needed to render the search HTML page
type SearchData struct {
	PageData
	Results []SearchDomain
	Prev    bool
	Next    bool
	PrevURL string
	NextURL string
}

type SearchDomain struct {
	Host string
	URL  string
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
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/search/1/%s", url.QueryEscape(search)))
		return
	} else {
		search = c.Param("query")
		if i, err := strconv.Atoi(c.Param("page")); err == nil {
			page = i
		}
	}

	var results []models.Domain

	searchFilter := fmt.Sprintf("%s%s%s", "%", search, "%")
	search2 := fmt.Sprintf("%s%s", "%", search)
	search4 := fmt.Sprintf("%s%s", search, "%")

	res := controller.db.
		Raw(fmt.Sprintf("SELECT * FROM domains WHERE host LIKE ? ORDER BY LENGTH(host), CASE WHEN host LIKE ? THEN 1 WHEN host LIKE ? THEN 2 WHEN host LIKE ? THEN 4 ELSE 3 END LIMIT %d OFFSET %d", resultsPerPage, resultsPerPage*(page-1)), searchFilter, search, search2, search4).
		Find(&results)

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
		pd.Results = append(pd.Results, SearchDomain{
			Host: host,
			URL:  fmt.Sprintf("/domains/%s", results[i].Host),
		})
	}

	if len(pd.Results) >= resultsPerPage {
		pd.Next = true
		pd.NextURL = fmt.Sprintf("/search/%d/%s", page+1, url.QueryEscape(search))
	}
	if page > 1 {
		pd.Prev = true
		pd.PrevURL = fmt.Sprintf("/search/%d/%s", page-1, url.QueryEscape(search))
	}

	c.HTML(http.StatusOK, "search.html", pd)
}
