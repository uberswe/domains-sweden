package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller Controller) Whois(c *gin.Context) {
	dpd := controller.DefaultPageData(c)
	dpd.Title = dpd.Trans("WHOIS Lookup")
	c.HTML(http.StatusOK, "whois.html", dpd)
}
