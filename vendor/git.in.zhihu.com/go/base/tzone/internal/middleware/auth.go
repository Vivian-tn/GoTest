package middleware

import (
	"net/http"
	"os"
)

func Auth(next http.Handler) http.Handler {
	enableAuth := os.Getenv("ZAE_SETTING_ENABLE_RPC_AUTH") == "1"
	expectedToken := os.Getenv("ZAE_SEC_APP_TOKEN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enableAuth {
			if r.Header.Get("X-ZONE-ORIGIN-TOKEN") != expectedToken &&
				r.Header.Get("X-ZONE-TARGET-TOKEN") != expectedToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
