package auth

// HMAC-SHA256 signing for HTTP requests
// Taken from https://learn.microsoft.com/en-us/azure/azure-app-configuration/rest-api-authentication-hmac#golang
// Modified sightly to be flexible with the signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SignRequestHMAC signs a HTTP request with HMAC-SHA256
func SignRequestHMAC(secret string, req *http.Request) error {
	method := req.Method
	host := req.URL.Host
	pathAndQuery := req.URL.Path

	if req.URL.RawQuery != "" {
		pathAndQuery = pathAndQuery + "?" + req.URL.RawQuery
	}

	content := []byte{}

	var err error

	if req.Body != nil {
		content, err = io.ReadAll(req.Body)
		if err != nil {
			// return err
			content = []byte{}
		}
	}

	req.Body = io.NopCloser(bytes.NewBuffer(content))

	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(http.TimeFormat)
	contentHash := GetContentHashBase64(content)
	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s", strings.ToUpper(method), pathAndQuery, timestamp, host, contentHash)
	signature := GetHmac(stringToSign, key)

	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("x-ms-date", timestamp)

	req.Header.Set("Authorization", "HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature="+signature)

	return nil
}

// Hash content with SHA256 and return the hash in base64
func GetContentHashBase64(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// Hash content with HMAC-SHA256 and return the hash in base64
func GetHmac(content string, key []byte) string {
	hmac := hmac.New(sha256.New, key)
	hmac.Write([]byte(content))

	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}
