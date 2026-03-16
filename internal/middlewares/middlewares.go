package middlewares

import (
	"compress/gzip"
	"context"
	"net/http"

	"github.com/flash1nho/GophKeeper/internal/authenticator"
)

type HTTPProvider struct {
	w http.ResponseWriter
	r *http.Request
}

func Decompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка при распаковке gzip", http.StatusBadRequest)
				return
			}

			defer gzReader.Close()

			r.Body = gzReader
		}

		next.ServeHTTP(w, r)
	})
}

func (p *HTTPProvider) GetCookie(_ context.Context, cookieName string) (string, error) {
	cookie, err := p.r.Cookie(cookieName)

	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (p *HTTPProvider) SetCookie(_ context.Context, cookieName, cookieValue string) error {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24 * 7,
	}

	http.SetCookie(p.w, cookie)

	return nil
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := authenticator.NewAuthenticator()
		ctx, err := auth.Authenticate(r.Context(), &HTTPProvider{w, r})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.Clone(ctx))
	})
}
