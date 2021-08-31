package server

import (
	"dap2pnet/rendezvous/middlewares"
	"dap2pnet/rendezvous/rendezvous"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	TLSCertPath string
	TLSKeytPath string
}

func Run(ren *rendezvous.Rendezvous) error {

	servConfig := &ServerConfig{
		TLSCertPath: "./certs/rendezvous.dap2p.net.pem",
		TLSKeytPath: "./certs/rendezvous.dap2p.net.key",
	}

	return InitializeEndpoints(servConfig, ren)

}

func InitializeEndpoints(servConfig *ServerConfig, ren *rendezvous.Rendezvous) error {
	gin.ForceConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(middlewares.Logger))

	peersGroup := router.Group("/peers/")
	peersGroup.Use(middlewares.SetPeerIdentity())

	InitPeerEndpoints(peersGroup, ren)

	return router.RunTLS(":6667", servConfig.TLSCertPath, servConfig.TLSKeytPath)
}
