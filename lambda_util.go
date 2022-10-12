package basic

import (
	"net/http"
	"os"
)

func getSchemeFromHeader(header http.Header) string {
	if scheme := header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}

func isStringContentTypeFromHeader(header http.Header) bool {
	contentType := header.Get(HeaderContentType)
	if len(contentType) == 0 {
		return false
	}

	return stringContentTypeTable[contentType]
}

const (
	envLambdaServerPort = "_LAMBDA_SERVER_PORT"
	envLambdaRuntimeAPI = "AWS_LAMBDA_RUNTIME_API"
)

func IsLambdaRuntime() bool {
	return isLambdaRuntime()
}

func isLambdaRuntime() bool {
	return os.Getenv(envLambdaServerPort) != "" || os.Getenv(envLambdaRuntimeAPI) != ""
}
