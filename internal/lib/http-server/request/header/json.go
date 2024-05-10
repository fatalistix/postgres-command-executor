package header

import (
	"net/http"
	"strings"
)

func HasApplicationJson(r *http.Request) bool {
	requestedHeader := r.Header.Get("Content-Type")
	requestedHeaderLower := strings.ToLower(requestedHeader)
	return strings.Contains(requestedHeaderLower, "application/json")
}
