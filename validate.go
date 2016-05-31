package securerequest

import (
	"net/http"
	"strconv"
	"time"
)

const defaultPeriod = 15 * time.Second

// Validate validates if a request has the appropriate authentication information
// It ignore requests othen than GET, POST, PUT, PATCH and DELETE
func Validate(r *http.Request, appSecrets map[string]string) bool {
	switch r.Method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
	default:
		return false
	}

	appKey := r.Header.Get(authAppHeader)
	if appKey == "" {
		return false
	}

	secret, ok := appSecrets[appKey]
	if !ok {
		return false
	}

	timestamp := r.Header.Get(authTimestampHeader)
	if timestamp == "" {
		return false
	}
	timeV, err := strconv.Atoi(timestamp)
	if err != nil {
		return false
	}
	t := time.Unix(0, int64(timeV*1000)) // micro to nano

	if time.Now().Sub(t) > defaultPeriod {
		return false
	}
	signature := r.Header.Get(authSignatureHeader)

	return signature == generateToken(appKey, secret, t, r.Method, r.URL.Path, r.Form, r.ContentLength)
}
