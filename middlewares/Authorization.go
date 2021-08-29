package middlewares

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AuthorizationMiddlewareErrIdentityNotFound = errors.New("identity not found. Proxy didn't send CN field")
)

func SetPeerIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		tlsHeader := c.GetHeader("X-Forwarded-Tls-Client-Cert-Info")
		if tlsHeader == "" {
			c.AbortWithError(http.StatusForbidden, AuthorizationMiddlewareErrIdentityNotFound)
		}

		q, err := url.QueryUnescape(tlsHeader)
		if err != nil {
			c.AbortWithError(http.StatusForbidden, AuthorizationMiddlewareErrIdentityNotFound)
		}

		q = strings.Split(q, ",")[0]
		identity := strings.ReplaceAll(strings.Split(q, "CN=")[1], "\"", "")
		c.Request.Header.Add("Peer-Identity", identity) // For logging purposes
		c.Set("Identity", identity)
	}
}
