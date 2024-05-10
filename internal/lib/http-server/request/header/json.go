package header

import (
	"net/http"
	"strings"
)

func HasApplicationJson(r *http.Request) bool {
	requestedHeader := r.Header.Get("Content-Type")
	return strings.ToLower(requestedHeader) == "application/json"
}
