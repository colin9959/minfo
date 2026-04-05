// Package middleware 提供日志和基础认证中间件。

package middleware

import (
	"crypto/subtle"
	"log"
	"net/http"
	"time"

	"minfo/internal/config"
	"minfo/internal/httpapi/transport"
)

// Logging 记录每个 HTTP 请求的方法、路径和耗时，然后继续执行下一个 Handler。
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// Authenticate 在配置了 Web 基础认证后校验请求凭据；未配置密码时直接放行。
func Authenticate(next http.Handler) http.Handler {
	username := config.Getenv("WEB_USERNAME", "")
	password := config.Getenv("WEB_PASSWORD", "")
	if password == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !matchesCredential(password, pass) || (username != "" && !matchesCredential(username, user)) {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"minfo\"")
			transport.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// matchesCredential 会判断认证信息是否匹配当前校验规则。
func matchesCredential(expected, actual string) bool {
	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}
