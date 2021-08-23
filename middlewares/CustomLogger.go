package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(param gin.LogFormatterParams) string {

	return fmt.Sprintf("%s -> tlsproxy -> %s - CN:%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.Request.Header.Get("X-Forwarded-For"),
		param.Request.RemoteAddr,
		param.Request.Header.Get("Authorization"),
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
