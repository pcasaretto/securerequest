package securerequest

import (
	"net/http"
	"strconv"
	"time"
)

// Sign takes a request, the app name and app secret and adds
// the appropriate headers to the request
// The request should already have its Content-Lenght set
func Sign(r *http.Request, app, secret string) {
	now := time.Now()
	timestamp := uint64(now.UnixNano() / int64(time.Millisecond))
	token := generateToken(app, secret, time.Now(), r.Method, r.URL.Path, r.URL.Query(), r.ContentLength)

	r.Header.Add(authAppHeader, app)
	r.Header.Add(authTimestampHeader, strconv.FormatUint(timestamp, 10))
	r.Header.Add(authSignatureHeader, token)
}
