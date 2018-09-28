package itbit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Config is the configuration to use to access the itBit API
type Config struct {
	// ClientSecret is the private API key. If not set, only public API endpoints will be available.
	ClientSecret string
	// APIBaseURL is the base url to use when accessing the API. If not set, the production API will be used.
	APIBaseURL string
	// ClientKey is the public client key
	ClientKey string
	// UserID is the user ID on the account to use
	UserID string
	nonce  int
}

func doReq(conf *Config, method, path string, signed bool, postData interface{}) (io.ReadCloser, error) {

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
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "reading %d error body", resp.StatusCode)
		}
		dec := json.NewDecoder(bytes.NewBuffer(body))
		apiErr := &itbitError{}
		if err := dec.Decode(apiErr); err != nil {
			return nil, errors.Wrapf(err, "parsing %d error body: %s", resp.StatusCode, body)
		}
		return nil, apiErr
	}
	return resp.Body, nil
}

func signRequest(conf *Config, r *http.Request, body []byte) error {
	debug("signing request")
	if conf.ClientSecret == "" || conf.ClientKey == "" {
		return errors.New("conf insufficient to make signed requests")
	}
	timestamp := time.Now().UnixNano() / 1000000
	conf.nonce++
	nonce := conf.nonce
	//urlStr := fmt.Sprintf("%s://%s%s", r.URL.Scheme, r.URL.Host, r.URL.Path)
	urlStr := r.URL.String()

	sigParts := []string{strings.ToUpper(r.Method), urlStr, string(body), fmt.Sprint(nonce), fmt.Sprint(timestamp)}
	debugf("signature parts: %#v", sigParts)
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(sigParts); err != nil {
		return errors.Wrap(err, "marshaling signature parts")
	}
	// take off the last byte, as the encoder adds newlines
	toHash := append([]byte(fmt.Sprintf("%d", nonce)), b.Bytes()[:len(b.Bytes())-1]...)
	debugf("to hash: %s", toHash)

	sha := sha256.New()
	if _, err := sha.Write(toHash); err != nil {
		return errors.Wrap(err, "writing to sha256 hasher")
	}

	shasum := sha.Sum(nil)
	debugf("sha256 sum: %x", shasum)

	hasher := hmac.New(sha512.New, []byte(conf.ClientSecret))
	toSign := append([]byte(urlStr), shasum...)

	debugf("to sign: %x", toSign)
	if _, err := hasher.Write(toSign); err != nil {
		return errors.Wrap(err, "writing url to sha512 hmac")
	}
	sigbytes := hasher.Sum(nil)

	debugf("sig bytes: %x", sigbytes)

	sig := base64.StdEncoding.EncodeToString(sigbytes)
	debugf("sig (hmac sha512): %s", sig)

	r.Header.Set("Authorization", fmt.Sprintf("%s:%s", conf.ClientKey, sig))
	r.Header.Set("x-auth-timestamp", fmt.Sprint(timestamp))
	r.Header.Set("x-auth-nonce", fmt.Sprint(nonce))

	debug("outgoing headers")
	for k, v := range r.Header {
		debug(k, v)
	}

	return nil

}

func makeURL(conf *Config, path string) string {
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
