package handlers

import (
	"net/http"
	"net/url"
	"strings"
)

func safeRedirect(s string) string {
	if s == "" {
		return ""
	}

	// accept only absolute-path, not protocol-relative
	if !strings.HasPrefix(s, "/") || strings.HasPrefix(s, "//") {
		return ""
	}

	return s
}

func redirectFromReferer(r *http.Request) string {
	ref := r.Header.Get("Referer")
	if ref == "" {
		return ""
	}

	u, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	// same-origin check
	if u.Host != r.Host {
		return ""
	}

	// avoid redirecting back to this handler or its subpaths
	if strings.HasPrefix(u.Path, r.URL.Path) {
		return ""
	}

	return safeRedirect(u.RequestURI())
}
