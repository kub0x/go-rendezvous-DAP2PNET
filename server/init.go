package server

import (
	"dap2pnet/rendezvous/middlewares"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	TLSCertPath string
	TLSKeytPath string
}

func Run() error {

	servConfig := &ServerConfig{
		TLSCertPath: "./certs/rendezvous.dap2p.net.pem",
		TLSKeytPath: "./certs/rendezvous.dap2p.net.key",
	}

	return InitializeEndpoints(servConfig)

}

func InitializeEndpoints(servConfig *ServerConfig) error {
	gin.ForceConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(middlewares.Logger))
	// Allow our tls proxy to pass its IP address
	router.TrustedProxies = []string{"127.0.0.1"}
	router.ForwardedByClientIP = true

	InitPeerEndpoints(router.Group("/peers/"))

	return router.RunTLS(":6668", servConfig.TLSCertPath, servConfig.TLSKeytPath)
}
