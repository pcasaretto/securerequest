package securerequest

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type example struct {
	Secrets         map[string]string
	Method          string
	AppHeader       string
	TimestampHeader string
	SignatureHeader string
	Body            io.Reader
	Expected        bool
}

var timestamp = time.Date(2016, time.May, 1, 0, 0, 0, 0, time.UTC)

func TestValidate(t *testing.T) {
	timestampGenerator = func() time.Time {
		return timestamp
	}
	tokenGenerator = func(appKey, appSecret string, timestamp time.Time, method string, path string, query url.Values, bodyLength int64) string {
		return "dummytoken"
	}

	secrets := map[string]string{
		"app_key": "super_sekret",
	}

	tests := []example{
		{nil, "BOGUS", "", "", "", nil, false},
		{secrets, "BOGUS", "", "", "", nil, true},
		{secrets, "GET", "", "", "", nil, false},
	}

	addTableTests(&tests, example{secrets, "GET", "", "", "", nil, false})

	for i, test := range tests {
		r, _ := http.NewRequest(test.Method, "", test.Body)
		headers := map[string]string{
			AuthAppHeader:       test.AppHeader,
			AuthTimestampHeader: test.TimestampHeader,
			AuthSignatureHeader: test.SignatureHeader,
		}
		for k, v := range headers {
			r.Header.Add(k, v)
		}
		if actual, expected := Validate(r, test.Secrets), test.Expected; actual != expected {
			t.Errorf("test %d:\n\texpected: %t\n\tgot: %t\nParameters: %+v", i, expected, actual, test)
		}
	}
}

func addTableTests(tests *[]example, template example) {
	methods := map[string]bool{
		"GET":     false,
		"POST":    false,
		"PUT":     false,
		"PATCH":   false,
		"DELETE":  false,
		"HEAD":    true,
		"OPTIONS": true,
	}
	appHeaderOptions := []struct {
		Header string
		Result bool
	}{
		{"", false},
		{"other_key", false},
		{"app_key", true},
	}
	timestampHeaderOptions := []struct {
		Header string
		Result bool
	}{
		{"", false},
		{"bogus", false},
		{timestampHeader(timestamp.Add(20 * time.Second)), false},
		{timestampHeader(timestamp.Add(-20 * time.Second)), false},
		{timestampHeader(timestamp.Add(15 * time.Second)), true},
		{timestampHeader(timestamp.Add(time.Second)), true},
		{timestampHeader(timestamp), true},
	}
	signatureHeaderOptions := []struct {
		Header string
		Result bool
	}{
		{"", false},
		{"bogus", false},
		{"dummytoken", true},
	}

	for method, result := range methods {
		for _, appHeader := range appHeaderOptions {
			for _, timestampHeader := range timestampHeaderOptions {
				for _, signatureHeader := range signatureHeaderOptions {
					template.Method = method
					template.AppHeader = appHeader.Header
					template.SignatureHeader = signatureHeader.Header
					template.TimestampHeader = timestampHeader.Header
					template.Expected = result || And(appHeader.Result, timestampHeader.Result, signatureHeader.Result)
					*tests = append(*tests, template)
				}
			}
		}
	}
}

func And(bools ...bool) bool {
	for _, b := range bools {
		if !b {
			return false
		}
	}
	return true
}

func timestampHeader(t time.Time) string {
	unixTimestamp := t.UnixNano() / int64(time.Millisecond)
	s := strconv.FormatInt(unixTimestamp, 10)
	return s
}
