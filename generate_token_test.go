package securerequest

import (
	"net/url"
	"testing"
	"time"
)

func TestNormalizeToken(t *testing.T) {
	timestamp := time.Date(2016, time.May, 1, 0, 0, 0, 0, time.UTC)
	query, _ := url.ParseQuery("z_param=z_value&m_param=m_value&a_param=a_value")
	tests := []struct {
		AppKey      string
		Secret      string
		Timestamp   time.Time
		Method      string
		Path        string
		QueryString url.Values
		BodyLength  int64
		Output      string
	}{
		{"kool_key", "super_sekret", timestamp, "POST", "/", query, 15, "app=kool_key&secret=kool_key&method=POST&path=/&t=1462060800000&a_param=a_value&m_param=m_value&z_param=z_value&body_length=15"},
		{"kool_key", "super_sekret", timestamp, "POST", "/some_route", query, 15, "app=kool_key&secret=kool_key&method=POST&path=/some_route&t=1462060800000&a_param=a_value&m_param=m_value&z_param=z_value&body_length=15"},
	}
	for i, test := range tests {
		actual := string(normalizeToken(test.AppKey, test.AppKey, test.Timestamp, test.Method, test.Path, test.QueryString, test.BodyLength))
		expected := test.Output
		if actual != expected {
			t.Errorf("test %d:\n\texpected: %s\n\tgot: %s", i, expected, actual)
		}
	}
}

func TestGenerateToken(t *testing.T) {
	timestamp := time.Date(2016, time.May, 1, 0, 0, 0, 0, time.UTC)
	query, _ := url.ParseQuery("z_param=z_value&m_param=m_value&a_param=a_value")
	tests := []struct {
		AppKey      string
		Secret      string
		Timestamp   time.Time
		Method      string
		Path        string
		QueryString url.Values
		BodyLength  int64
		Output      string
	}{
		//TODO: validate hashes
		{"kool_key", "super_sekret", timestamp, "POST", "/", query, 15, "bMlnMRQ8epVfLzkdR76fCL/QxYc1yz5kX2tL/6r0BCI="},
		{"kool_key", "super_sekret", timestamp, "POST", "/some_route", query, 15, "pLF3Apie1dv8Wv8AMdLSKw2w3SNLMx3ARmqWR2MLowI="},
	}
	for i, test := range tests {
		actual := generateToken(test.AppKey, test.AppKey, test.Timestamp, test.Method, test.Path, test.QueryString, test.BodyLength)
		expected := test.Output
		if actual != expected {
			t.Errorf("test %d:\n\texpected: %s\n\tgot: %s", i, expected, actual)
		}
	}
}
