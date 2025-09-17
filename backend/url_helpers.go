package main

import (
	"net/http"
	"net/url"
)

func buildCallbackURL(r *http.Request, ticketID string) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}

	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   "/api/irma/callback",
	}
	q := u.Query()
	q.Set("ticketId", ticketID)
	u.RawQuery = q.Encode()
	return u.String()
}
