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
	if appSecrets == nil {
		return false
	}
	switch r.Method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
	default:
		return true
	}

	appKey := r.Header.Get(AuthAppHeader)
	if appKey == "" {
		return false
	}

	secret, ok := appSecrets[appKey]
	if !ok {
		return false
	}

	timestamp := r.Header.Get(AuthTimestampHeader)
	if timestamp == "" {
		return false
	}
	timeV, err := strconv.Atoi(timestamp)
	if err != nil {
		return false
	}
	t := time.Unix(int64(timeV/1000), int64(timeV%1000))

	diff := timestampGenerator().Sub(t)
	if diff > defaultPeriod || diff < -defaultPeriod {
		return false
	}
	signature := r.Header.Get(AuthSignatureHeader)

	return signature == tokenGenerator(appKey, secret, t, r.Method, r.URL.Path, r.URL.Query(), r.ContentLength)
}
