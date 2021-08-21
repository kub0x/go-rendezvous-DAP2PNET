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
	router.LoadHTMLGlob("templates/*")
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(middlewares.Logger))

	println("All good!")

	if gin.IsDebugging() {

	}

	return router.RunTLS(":6668", servConfig.TLSCertPath, servConfig.TLSKeytPath)
}