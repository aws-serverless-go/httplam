package basic

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"strconv"
	"strings"
)

func NewHTTPRequest(ctx context.Context, from *events.APIGatewayV2HTTPRequest) (req *http.Request, err error) {
	rawProtocol := from.RequestContext.HTTP.Protocol
	protoDiv := strings.Index(rawProtocol, "/")
	if protoDiv == -1 {
		return
	}

	proto := rawProtocol[protoDiv+1:]
	protoDiv = strings.Index(proto, ".")
	if protoDiv == -1 {
		return
	}

	major, err := strconv.Atoi(proto[:protoDiv])
	if err != nil {
		return
	}
	minor, err := strconv.Atoi(proto[protoDiv+1:])
	if err != nil {
		return
	}

	var body []byte
	if from.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(from.Body)
		if err != nil {
			return
		}
	} else {
		body = []byte(from.Body)
	}

	header := make(http.Header)
	for k, v := range from.Headers {
		header.Set(k, v)
	}

	var rawURL strings.Builder
	rawURL.WriteString(getSchemeFromHeader(header))
	rawURL.WriteString("://")
	rawURL.WriteString(header.Get(HeaderHost))
	rawURL.WriteString(from.RawPath)
	if len(from.RawQueryString) > 0 {
		rawURL.WriteRune('?')
		rawURL.WriteString(from.RawQueryString)
	}

	req, err = http.NewRequestWithContext(ctx, from.RequestContext.HTTP.Method, rawURL.String(), bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Proto = rawProtocol
	req.ProtoMajor = major
	req.ProtoMinor = minor
	req.RemoteAddr = from.RequestContext.HTTP.SourceIP
	req.Header = header
	return
}
