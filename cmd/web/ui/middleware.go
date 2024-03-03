package ui

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/briancbarrow/gitfit-go/cmd/web"
)

func generateNonce() (string, error) {
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}
	return base64.StdEncoding.EncodeToString(nonce), nil
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonceString, err := generateNonce()
		if err != nil {
			fmt.Println("error generating nonce")
		}
		ctx := context.WithValue(r.Context(), web.NonceValue, nonceString)
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self' 'unsafe-eval' 'nonce-"+nonceString+"'")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
