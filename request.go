package itbit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config is the configuration to use to access the itBit API
type Config struct {
	// APIKey is the private API key. If not set, only public API endpoints will be available.
	APIKey string
	// APIBaseURL is the base url to use when accessing the API. If not set, the production API will be used.
	APIBaseURL string
	// ClientKey is the public client key
	ClientKey string
}

func doReq(conf Config, method, path string, signed bool, postData interface{}) (io.ReadCloser, error) {

	path = makeURL(conf, path)

	var (
		body io.Reader
		data []byte
	)
	if postData != nil {

		var err error
		data, err = json.Marshal(postData)
		if err != nil {
			return nil, errors.Wrap(err, "marshaling json")
		}
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(strings.ToUpper(method), path, body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	req.Header.Set("content-type", "application/json")

	if signed {
		if err := signRequest(conf, req, data); err != nil {
			return nil, errors.Wrap(err, "signing request")
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("%d error", resp.StatusCode)
	}
	return resp.Body, nil
}

func signRequest(conf Config, r *http.Request, body []byte) error {
	if conf.APIKey == "" || conf.ClientKey == "" {
		return errors.New("conf insufficient to make signed requests")
	}
	timestamp := time.Now().UnixNano() / 1000
	nonce := timestamp - 42
	sigParts := []string{strings.ToUpper(r.Method), r.URL.String(), string(body), fmt.Sprint(nonce), fmt.Sprint(timestamp)}
	toSign, err := json.Marshal(sigParts)
	if err != nil {
		return errors.Wrap(err, "marshaling signature parts")
	}

	sha := sha256.New()
	if _, err := sha.Write([]byte(fmt.Sprint(nonce))); err != nil {
		return errors.Wrap(err, "writing to sha256 hasher")
	}
	sha.Write(toSign)

	hasher := hmac.New(sha512.New, []byte(conf.APIKey))

	if _, err := hasher.Write(sha.Sum(nil)); err != nil {
		return errors.Wrap(err, "writing to sha512 hmac")
	}
	sig := hasher.Sum(nil)

	r.Header.Set("authorization", fmt.Sprintf("%s:%X", conf.ClientKey, sig))
	r.Header.Set("x-auth-timestamp", fmt.Sprint(timestamp))
	r.Header.Set("x-auth-nonce", fmt.Sprint(nonce))

	return nil

}

func makeURL(conf Config, path string) string {
	base := conf.APIBaseURL
	if base == "" {
		base = `https://api.itbit.com/v1`
	}
	if path == "" {
		return base
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return base + path
}
