package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnSuscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("The peer is: " + c.GetHeader("Authorization"))
		payload, _ := c.GetRawData()
		log.Println("The data is: " + string(payload))
		c.Data(http.StatusForbidden, "text/html", []byte("Hello"))
	}
}
