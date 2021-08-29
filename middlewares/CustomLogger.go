package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(param gin.LogFormatterParams) string {

	return fmt.Sprintf("client: %s -> traefik: %s -> server - CN:%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.Request.Header.Get("X-Forwarded-For"),
		param.ClientIP,
		param.Request.Header.Get("Peer-Identity"),
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)

}
