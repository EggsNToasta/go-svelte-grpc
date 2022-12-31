package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mpiorowski/golang"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	ENV           = utils.MustGetenv("ENV")
	PORT          = utils.MustGetenv("PORT")
	EMAIL_API_KEY = utils.MustGetenv("EMAIL_API_KEY")
	CLIENT_URL    = utils.MustGetenv("CLIENT_URL")
	EMAIL_FROM    = utils.MustGetenv("EMAIL_FROM")
	EMAIL_NAME    = utils.MustGetenv("EMAIL_NAME")
)

func main() {
	if ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.POST("/pubsub", getPubSubEmail)
	err := router.Run(fmt.Sprintf("0.0.0.0:%v", PORT))
	if err != nil {
		panic(err)
	}
}

func getPubSubEmail(c *gin.Context) {

	message, err := utils.SubscribePubSub(c)
	if err != nil {
		log.Printf("utils.SubscribePubSub: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var email utils.Email
	err = json.Unmarshal(message.Message.Data, &email)
	if err != nil {
		log.Printf("json.Unmarshal: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from := mail.NewEmail(EMAIL_NAME, EMAIL_FROM)
	to := mail.NewEmail(email.To, email.To)

	subject, body := getTemplate(email.Template, email.Html)

	msg := mail.NewSingleEmail(from, subject, to, "", body)
	client := sendgrid.NewSendClient(EMAIL_API_KEY)
	response, err := client.Send(msg)
	if err != nil {
		log.Printf("client.Send: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Email send: %v", response.Headers)
	c.JSON(http.StatusOK, gin.H{"message": "email send"})
}