package apollo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	httpHeaderAuthorization = "Authorization"
	httpHeaderTimestamp     = "Timestamp"

	authorizationFormat = "Apollo %s:%s"

	delimiter = "\n"
	question  = "?"
)

//apollo 授权
type AuthSignature struct {
}

// HTTPHeaders
func (t *AuthSignature) HTTPHeaders(url string, appID string, secret string) map[string][]string {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := strconv.FormatInt(ms, 10)
	pathWithQuery := url2PathWithQuery(url)

	stringToSign := timestamp + delimiter + pathWithQuery
	signature := signString(stringToSign, secret)
	headers := make(map[string][]string, 2)

	signatures := make([]string, 0, 1)
	signatures = append(signatures, fmt.Sprintf(authorizationFormat, appID, signature))
	headers[httpHeaderAuthorization] = signatures

	timestamps := make([]string, 0, 1)
	timestamps = append(timestamps, timestamp)
	headers[httpHeaderTimestamp] = timestamps
	return headers
}

func signString(stringToSign string, accessKeySecret string) string {
	key := []byte(accessKeySecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func url2PathWithQuery(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	pathWithQuery := u.Path
	if len(u.RawQuery) > 0 {
		pathWithQuery += question + u.RawQuery
	}
	return pathWithQuery
}

