package server

import (
	"dap2pnet/rendezvous/handlers"

	"dap2pnet/rendezvous/rendezvous"

	"github.com/gin-gonic/gin"
)

func InitPeerEndpoints(router *gin.RouterGroup, ren *rendezvous.Rendezvous) {
	router.POST("/subscribe", handlers.OnSubscribe(ren))
	router.GET("/", handlers.OnGetPeers(ren))
}
