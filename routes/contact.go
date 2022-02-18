package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	email2 "github.com/uberswe/domains-sweden/email"
	"log"
	"net/http"
)

// Contact ...
func (controller Controller) Contact(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Contact Us")
	c.HTML(http.StatusOK, "contact.html", pd)
}

// ContactPost ...
func (controller Controller) ContactPost(c *gin.Context) {
	var err error
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Contact Us")

	email := c.PostForm("email")
	name := c.PostForm("name")
	content := c.PostForm("content")

	validate := validator.New()
	err = validate.Var(email, "required,email")

	if len(name) < 1 {
		err = errors.New("name is not set")
	}

	if len(content) < 10 {
		err = errors.New("content is not set")
	}

	if err != nil {
		log.Println(err)
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: pd.Trans("Please make sure you include a name, email and message"),
		})
		c.HTML(http.StatusBadRequest, "contact.html", pd)
		return
	}

	go controller.sendContactEmail(name, email, content, pd.Trans)

	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pd.Trans("Message sent!"),
	})

	c.HTML(http.StatusOK, "contact.html", pd)
}

func (controller Controller) sendContactEmail(name string, email string, content string, trans func(string) string) {
	emailService := email2.New(controller.config)
	emailService.Send("contact@domaner.xyz", trans("[Domäner.xyz] Contact Form"), fmt.Sprintf(trans("A new email was sent from Domäner.xyz.\n\nName: %s\nEmail: %s\nContent: %s\n"), name, email, content))
}
