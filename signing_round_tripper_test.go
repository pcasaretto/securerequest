package securerequest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pcasaretto/securerequest"
)

func TestNewSigningRoundTripper(t *testing.T) {
	secrets := map[string]string{"app": "secret"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
		if !securerequest.Validate(r, secrets) {
			t.Error("request was not valid")
		}
	}))
	defer ts.Close()

	tr := securerequest.NewSigningRoundTripper(nil, "app", "secret")
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{Transport: tr}
	client.Do(req)
}
