package middleware

import (
	"context"
	"net/http"
)

type ctxKey string

const (
	callerKey ctxKey = "caller-ctx-key"
)

type Caller struct {
	App     string
	Service string
	Region  string
}

func InjectCaller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		caller := &Caller{
			App:     r.Header.Get("X-ZONE-ORIGIN-APP"),
			Service: r.Header.Get("X-ZONE-ORIGIN-UNIT"),
			Region:  r.Header.Get("X-ZONE-ORIGIN-REGION"),
		}
		ctx := context.WithValue(r.Context(), callerKey, caller)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CallerFromContext(ctx context.Context) *Caller {
	if v := ctx.Value(callerKey); v != nil {
		return v.(*Caller)
	}
	return nil
}
