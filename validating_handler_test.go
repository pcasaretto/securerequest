package securerequest_test

import (
	"net/http"
	"testing"

	"github.com/pcasaretto/securerequest"
)

func TestNewValidatingHandler(t *testing.T) {
	var h http.Handler
	var secrets map[string]string
	assertPanic(t, func() {
		securerequest.NewValidatingHandler(h, secrets)
	})

	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	assertPanic(t, func() {
		securerequest.NewValidatingHandler(h, secrets)
	})

	secrets = make(map[string]string)
	assertPanic(t, func() {
		securerequest.NewValidatingHandler(h, secrets)
	})

	secrets["app"] = "secret"

	assertNoPanic(t, func() {
		securerequest.NewValidatingHandler(h, secrets)
	})
}
