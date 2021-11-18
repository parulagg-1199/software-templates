package auth

import (
	"context"
	"net/http"
)

// CtxKey is used for saving token data in context correspondance to this key
type CtxKey struct{}

// Middleware is a middleware function to check authorization
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if auth != "" {
			// Write your fancy token introspection logic here and if valid user then pass appropriate key in header
			// IMPORTANT: DO NOT HANDLE UNAUTHORISED USER HERE

			ctx = context.WithValue(ctx, CtxKey{}, auth)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
