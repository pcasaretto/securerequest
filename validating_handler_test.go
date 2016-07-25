package securerequest_test

import (
	"net/http"
	"net/http/httptest"
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

func TestValidatingHandler_ServeHTTP(t *testing.T) {
	var (
		h       http.Handler
		secrets map[string]string
		w       *httptest.ResponseRecorder
		req     *http.Request
		v       *securerequest.ValidatingHandler
	)
	secrets = make(map[string]string)

	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	secrets["app"] = "secret"

	w = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	v = securerequest.NewValidatingHandler(h, secrets)
	v.ServeHTTP(w, req)
	if actual, expected := w.Code, http.StatusUnauthorized; actual != expected {
		t.Errorf("\texpected: %d\n\tgot: %d", expected, actual)
	}
	if actual, expected := w.Body.String(), securerequest.UnauthorizedMessage; actual != expected {
		t.Errorf("\texpected: %s\n\tgot: %s", expected, actual)
	}
	if actual, expected := w.Header().Get("Content-Type"), "plain-text"; actual != expected {
		t.Errorf("\texpected: %s\n\tgot: %s", expected, actual)
	}

	var called bool
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	securerequest.Sign(req, "app", "secret")
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		called = true
	})
	v = securerequest.NewValidatingHandler(h, secrets)
	v.ServeHTTP(w, req)
	if actual, expected := called, true; actual != expected {
		t.Errorf("\texpected: %b\n\tgot: %b", expected, actual)
	}
	if actual, expected := w.Code, http.StatusOK; actual != expected {
		t.Errorf("\texpected: %d\n\tgot: %d", expected, actual)
	}

}
