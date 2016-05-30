package securerequest

import "net/http"

func Sign(r *http.Request) {
}

func Validate(r *http.Request) bool {
	return false
}

type SigningRoundTripper struct {
	rt http.RoundTripper
}

func NewSigningRoundTripper(rt http.RoundTripper) *SigningRoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &SigningRoundTripper{rt}
}

func (srt *SigningRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	Sign(r)
	return srt.rt.RoundTrip(r)
}

type ValidatingHandler struct {
	handler http.Handler
}

func NewValidatingHandler(h http.Handler) *ValidatingHandler {
	return &ValidatingHandler{h}
}

const unauthorized_message = "Unauthorized access."

func (vh *ValidatingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !Validate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(unauthorized_message))
		return
	}
	vh.handler.ServeHTTP(w, r)
}
