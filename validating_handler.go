package securerequest

import "net/http"

// ValidatingHandler is a middleware that authenticates requests.
type ValidatingHandler struct {
	handler    http.Handler
	appSecrets map[string]string
}

// NewValidatingHandler takes a http.Handler and a app -> secret map and
// returns a ValidatingHandler that uses the given secrets to validate requests
func NewValidatingHandler(h http.Handler, appSecrets map[string]string) *ValidatingHandler {
	return &ValidatingHandler{h, appSecrets}
}

const unauthorizedMessage = "Unauthorized access."

// ServeHTTP validates requests. If a request passes validation it passes along
// to the next Handler. If it doesn't it will respond with 403 and plain-text
// message "Unauthorized access."
func (vh *ValidatingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !Validate(r, vh.appSecrets) {
		w.Header().Add("Content-Type", "plain-text")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(unauthorizedMessage))
		return
	}
	vh.handler.ServeHTTP(w, r)
}
