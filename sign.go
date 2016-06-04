package securerequest

import (
	"net/http"
	"strconv"
	"time"
)

var timestampGenerator = func() time.Time {
	return time.Now()
}
var tokenGenerator = generateToken

// Sign takes a request, the app name and app secret and adds
// the appropriate headers to the request
// The request should already have its Content-Lenght set
func Sign(r *http.Request, app, secret string) {
	now := timestampGenerator()
	timestamp := uint64(now.UnixNano() / int64(time.Millisecond))
	token := tokenGenerator(app, secret, time.Now(), r.Method, r.URL.Path, r.URL.Query(), r.ContentLength)

	r.Header.Add(AuthAppHeader, app)
	r.Header.Add(AuthTimestampHeader, strconv.FormatUint(timestamp, 10))
	r.Header.Add(AuthSignatureHeader, token)
}
