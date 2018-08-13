package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	xhttp "net/http"
	xurl "net/url"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
)

const OptionTypeUrlencode = "application/x-www-form-urlencoded"
const OptionTypeJson = "application/json"
const OptionTypeText = "text/plain"
const OptionTypeXml = "text/xml"

const OptionResponseTypeText = "text"
const OptionResponseTypeJson = "json"
const OptionResponseTypeByte = "byte"

type Options struct {
	Url           string
	Method        string
	Type          string
	ResponseType  string
	Data          interface{}
	Headers       map[string]string
	RedirectCount int
}

var ca *x509.CertPool

func init() {
	ca = x509.NewCertPool()
	ca.AppendCertsFromPEM(pemCerts)
}

func NewClient() *xhttp.Client {
	return &xhttp.Client{
		Transport: &xhttp.Transport{
			TLSClientConfig:   &tls.Config{RootCAs: ca},
			DisableKeepAlives: false,
		},
	}
}

func Send(options *Options) (interface{}, error) {

	client := &xhttp.Client{
		Transport: &xhttp.Transport{
			TLSClientConfig:   &tls.Config{RootCAs: ca},
			DisableKeepAlives: false,
		},
	}

	var url = options.Url
	var resp *xhttp.Response
	var req *xhttp.Request
	var err error

	if options.Method == "POST" {

		var body []byte = nil

		if strings.Contains(options.Type, "json") {

			body, err = json.Marshal(options.Data)

			if err != nil {
				return nil, err
			}

		} else if strings.Contains(options.Type, "text") {

			body = []byte(dynamic.StringValue(options.Data, ""))

		} else {

			idx := 0
			b := bytes.NewBuffer(nil)

			dynamic.Each(options.Data, func(key interface{}, value interface{}) bool {

				if idx != 0 {
					b.WriteString("&")
				}

				b.WriteString(dynamic.StringValue(key, ""))
				b.WriteString("=")
				b.WriteString(xurl.QueryEscape(dynamic.StringValue(value, "")))

				idx = idx + 1

				return true
			})

			body = b.Bytes()
		}

		req, err = xhttp.NewRequest("POST", url, bytes.NewReader(body))

		if err == nil {

			req.Header.Set("Content-Type", options.Type+"; charset=utf-8")

			if options.Headers != nil {
				for key, value := range options.Headers {
					req.Header.Set(key, value)
				}
			}

			resp, err = client.Do(req)
		}

	} else {

		idx := 0

		b := bytes.NewBuffer(nil)

		dynamic.Each(options.Data, func(key interface{}, value interface{}) bool {

			if idx != 0 {
				b.WriteString("&")
			}

			b.WriteString(dynamic.StringValue(key, ""))
			b.WriteString("=")
			b.WriteString(xurl.QueryEscape(dynamic.StringValue(value, "")))

			idx = idx + 1

			return true
		})

		idx = strings.Index(url, "?")

		if idx >= 0 {
			if idx+1 == len(url) {
				url = url + b.String()
			} else {
				url = url + "&" + b.String()
			}
		} else {
			url = url + "?" + b.String()
		}

		req, err = xhttp.NewRequest("GET", url, nil)

		if err == nil {
			if options.Headers != nil {
				for key, value := range options.Headers {
					req.Header.Add(key, value)
				}
			}
			resp, err = client.Do(req)
		}

	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {

		b := bytes.NewBuffer(nil)

		_, err = b.ReadFrom(resp.Body)

		resp.Body.Close()

		if err != nil && err != io.EOF {
			return nil, err
		}

		if options.ResponseType == "json" {
			var data interface{} = nil
			err := json.Unmarshal(b.Bytes(), &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		} else if options.ResponseType == "byte" {
			return b.Bytes(), nil
		} else {
			return string(b.Bytes()), nil
		}

	} else {

		if resp.StatusCode == 302 && options.RedirectCount > 0 {
			options.Url = resp.Header.Get("Location")
			options.RedirectCount = options.RedirectCount - 1
			fmt.Println("[KK] Redirect", options.Url)
			return Send(options)
		}

		b := bytes.NewBuffer(nil)

		_, err = b.ReadFrom(resp.Body)

		resp.Body.Close()

		if err != nil && err != io.EOF {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("[%d] %s", resp.StatusCode, string(b.Bytes())))
	}
}
