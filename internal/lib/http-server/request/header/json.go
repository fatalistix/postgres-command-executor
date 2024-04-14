package header

import "net/http"

func HasApplicationJson(r *http.Request) bool {
	requestedHeader := r.Header.Get("Content-Type")
	return requestedHeader == "application/json"
}
