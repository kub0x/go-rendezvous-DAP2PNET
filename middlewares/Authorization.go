package middlewares

import "github.com/gin-gonic/gin"

func SetPeerIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.GetHeader("Authorization")
		if identity == "" {

		}

		c.Set("Identity", identity)
	}
}
