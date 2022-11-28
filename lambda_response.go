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

func NewAPIGatewayV2HTTPResponseBuilder(target *events.APIGatewayV2HTTPResponse) (res *APIGatewayV2HTTPResponseBuilder) {
	res = &APIGatewayV2HTTPResponseBuilder{
		header:      make(http.Header),
		buildTarget: target,
	}

	if res.buildTarget == nil {
		res.buildTarget = new(events.APIGatewayV2HTTPResponse)
	}

	return
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

func (w *APIGatewayV2HTTPResponseBuilder) Build() (*events.APIGatewayV2HTTPResponse, error) {
	return w.build()
}

func (w *APIGatewayV2HTTPResponseBuilder) build() (res *events.APIGatewayV2HTTPResponse, err error) {
	res = w.buildTarget
	if w.done {
		return
	}

	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}

	body, err := ioutil.ReadAll(&w.buffer)
	if err != nil {
		res = nil
		return
	}

	if len(body) > 0 {
		if isStringContentTypeFromHeader(w.header) {
			res.Body = string(body)
		} else {
			res.IsBase64Encoded = true
			res.Body = base64.StdEncoding.EncodeToString(body)
		}
	}

	res.Headers = make(map[string]string)
	res.MultiValueHeaders = make(map[string][]string)

	for k, v := range w.header {
		// known-headers
		switch k {
		case HeaderSetCookie:
			w.setCookies(v...)
			continue
		}

		if len(v) == 1 {
			res.Headers[k] = v[0]
		} else {
			res.MultiValueHeaders[k] = v
		}
	}

	w.done = true
	return
}

func (w *APIGatewayV2HTTPResponseBuilder) IsDone() bool {
	return w.done
}
