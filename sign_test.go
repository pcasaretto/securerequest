package securerequest

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	timestamp := time.Date(2016, time.May, 1, 0, 0, 0, 0, time.UTC)
	timestampGenerator = func() time.Time {
		return timestamp
	}
	tokenGenerator = func(appKey, appSecret string, timestamp time.Time, method string, path string, query url.Values, bodyLength int64) string {
		return "dummytoken"
	}
	r, _ := http.NewRequest("GET", "http://secret.app", nil)
	Sign(r, "app_key", "super_sekret")

	tests := []struct {
		Header   string
		Expected string
	}{
		{AuthAppHeader, "app_key"},
		{AuthTimestampHeader, "1462060800000"},
		{AuthSignatureHeader, "dummytoken"},
	}

	for i, test := range tests {
		if actual, expected := r.Header.Get(test.Header), test.Expected; actual != expected {
			t.Errorf("test %d:\n\texpected: %s\n\tgot: %s", i, expected, actual)
		}
	}

}
