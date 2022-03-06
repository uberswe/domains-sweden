package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/domains-sweden/lang"
	"log"
	"net/http"
)

func (controller Controller) Language(c *gin.Context) {
	dpd := controller.DefaultPageData(c)

	l := c.Param("lang")

	session := sessions.Default(c)
	if l == "en" {
		session.Set(lang.CookieKey, "en")
	} else if l == "sv" {
		session.Set(lang.CookieKey, "sv")
	} else {
		c.HTML(http.StatusNotFound, "404.html", dpd)
		return
	}
	err := session.Save()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "404.html", dpd)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
	return
}
