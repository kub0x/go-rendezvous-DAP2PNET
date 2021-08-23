package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnSuscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, _ := c.GetRawData()
		c.Data(http.StatusForbidden, "text/html", payload)
	}
}
