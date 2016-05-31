package securerequest

import "net/http"

// A SigningRoundTripper implements http.RoundTripper
// by decorating another http.RoundTripper
// It signs signs all requests with the given app and secret string
// before passing them along
type SigningRoundTripper struct {
	rt     http.RoundTripper
	app    string
	secret string
}

// NewSigningRoundTripper returns a SigningRoundTripper that
// signs requests with the given app and secret
// if given a nil http.RoundTripper, it uses http.DefaultTransport
func NewSigningRoundTripper(rt http.RoundTripper, app, secret string) *SigningRoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &SigningRoundTripper{rt, app, secret}
}

// RoundTrip Signs the request before passing it along to the
// decorated http.RoundTripper
func (srt *SigningRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	Sign(r, srt.app, srt.secret)
	return srt.rt.RoundTrip(r)
}
