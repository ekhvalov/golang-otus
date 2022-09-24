package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func loggingMiddleware(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			// 66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
			remoteIP := strings.Split(r.RemoteAddr, ":")[0]
			logger.Info(fmt.Sprintf(
				"%s [%s] %s %s %s %d %v",
				remoteIP,
				start.Format("01/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.RequestURI,
				r.Proto,
				time.Since(start)/time.Microsecond,
				r.UserAgent()),
			)
		}()
		next.ServeHTTP(w, r)
	})
}
