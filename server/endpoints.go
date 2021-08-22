package server

import (
	"dap2pnet/rendezvous/handlers"

	"github.com/gin-gonic/gin"
)

func InitPeerEndpoints(router *gin.RouterGroup) {
	router.POST("/subscribe", handlers.OnSuscribe())
}
