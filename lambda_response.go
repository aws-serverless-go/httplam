package httplam

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"net"
	"net/http"
)

func NewAPIGatewayV2HTTPResponseBuilder(target *events.APIGatewayV2HTTPResponse) *APIGatewayV2HTTPResponseBuilder {
	return &APIGatewayV2HTTPResponseBuilder{
		header:      make(http.Header),
		buildTarget: target,
	}
}

var (
	_ http.ResponseWriter = (*APIGatewayV2HTTPResponseBuilder)(nil)
	_ http.Hijacker       = (*APIGatewayV2HTTPResponseBuilder)(nil)
)

type APIGatewayV2HTTPResponseBuilder struct {
	header      http.Header
	buffer      bytes.Buffer
	done        bool
	buildTarget *events.APIGatewayV2HTTPResponse
}

func (w *APIGatewayV2HTTPResponseBuilder) Header() http.Header {
	return w.header
}

func (w *APIGatewayV2HTTPResponseBuilder) Write(bytes []byte) (int, error) {
	return w.buffer.Write(bytes)
}

func (w *APIGatewayV2HTTPResponseBuilder) WriteHeader(statusCode int) {
	w.buildTarget.StatusCode = statusCode
}

func (w *APIGatewayV2HTTPResponseBuilder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, errors.New("not supported")
}

func (w *APIGatewayV2HTTPResponseBuilder) SetCookie(cookie *http.Cookie) {
	v := cookie.String()
	if len(v) == 0 {
		return
	}

	w.setCookies(v)
}

func (w *APIGatewayV2HTTPResponseBuilder) setCookies(v ...string) {
	w.buildTarget.Cookies = append(w.buildTarget.Cookies, v...)
}

func (w *APIGatewayV2HTTPResponseBuilder) Build() error {
	return w.build()
}

func (w *APIGatewayV2HTTPResponseBuilder) build() error {
	if w.done {
		return nil
	}

	if w.buildTarget.StatusCode == 0 {
		w.buildTarget.StatusCode = http.StatusOK
	}

	body, err := ioutil.ReadAll(&w.buffer)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		if isStringContentTypeFromHeader(w.header) {
			w.buildTarget.Body = string(body)
		} else {
			w.buildTarget.IsBase64Encoded = true
			w.buildTarget.Body = base64.StdEncoding.EncodeToString(body)
		}
	}

	w.buildTarget.Headers = make(map[string]string)
	w.buildTarget.MultiValueHeaders = make(map[string][]string)

	for k, v := range w.header {
		// known-headers
		switch k {
		case HeaderSetCookie:
			w.setCookies(v...)
			continue
		}

		if len(v) == 1 {
			w.buildTarget.Headers[k] = v[0]
		} else {
			w.buildTarget.MultiValueHeaders[k] = v
		}
	}

	w.done = true
	return nil
}

func (w *APIGatewayV2HTTPResponseBuilder) IsDone() bool {
	return w.done
}
