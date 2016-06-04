package securerequest

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

func generateToken(appKey, appSecret string, timestamp time.Time, method string, path string, query url.Values, bodyLength int64) string {

	token := normalizeToken(appKey, appSecret, timestamp, method, path, query, bodyLength)

	hasher := sha256.New()
	hasher.Write(token)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func normalizeToken(appKey, appSecret string, timestamp time.Time, method string, path string, query url.Values, bodyLength int64) []byte {
	var output bytes.Buffer
	unixTimestamp := timestamp.UnixNano() / int64(time.Millisecond)

	tokenString := fmt.Sprintf("app=%s&secret=%s&method=%s&path=%s&t=%d", appKey, appSecret, method, path, unixTimestamp)
	output.WriteString(tokenString)

	if len(query) > 0 {
		output.WriteRune('&')
		output.WriteString(query.Encode())
	}

	if bodyLength > 0 {
		output.WriteString(fmt.Sprintf("&body_length=%d", bodyLength))
	}
	return output.Bytes()
}
