package middlewares

import "github.com/gin-gonic/gin"

func SetPeerIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.GetHeader("Authorization")
		tlsHeader := c.GetHeader("X-Forwarded-Tls-Client-Cert-Info")
		println("TLS HEADER: " + tlsHeader)

		if identity == "" {

		}

		c.Set("Identity", identity)
	}
}
