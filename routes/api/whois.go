package api

import (
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/whois"
	"net/http"
	"strings"
)

type WhoisRequest struct {
	Domain string `json:"domain"`
}

type WhoisResponse struct {
	whois.Whois
}

func (controller Controller) Whois(c *gin.Context) {
	var request WhoisRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	if !strings.HasSuffix(request.Domain, ".se") && !strings.HasSuffix(request.Domain, ".nu") {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	whoisService := whois.New(controller.db)
	res, err := whoisService.Domain(request.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	var response WhoisResponse
	response.Raw = res.Raw
	response.Fetched = res.Fetched
	c.JSON(http.StatusOK, response)
}
