package middleware

import (
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// start := time.Now()
		// requestID := fmt.Sprintf("%016x", rand.Int())[:10]
		// newContext := context.WithValue(req.Context(), "request_id", requestID)

		// logger.Middleware().Info(newContext, logger.Fields{
		// 	"url":         req.URL,
		// 	"method":      req.Method,
		// 	"remote_addr": req.RemoteAddr,
		// 	// "server_status": req.,
		// })

		next.ServeHTTP(w, req)

		// logger.Middleware().Info(newContext, logger.Fields{
		// 	"work_time": time.Since(start),
		// })
	})
}
