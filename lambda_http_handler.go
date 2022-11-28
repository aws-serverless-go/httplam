package httplam

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var _ lambda.Handler = (*defaultHandler)(nil)

type defaultHandler struct {
	handler http.Handler
}

func (h *defaultHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	var request events.APIGatewayV2HTTPRequest
	err := json.Unmarshal(payload, &request)
	if err != nil {
		return nil, err
	}

	req, err := NewHTTPRequest(ctx, &request)
	if err != nil {
		return nil, err
	}

	var response events.APIGatewayV2HTTPResponse
	rw := NewAPIGatewayV2HTTPResponseBuilder(&response)

	h.handler.ServeHTTP(rw, req)
	_, err = rw.Build()
	if err != nil {
		return nil, err
	}

	return json.Marshal(response)
}

func StartLambdaWithAPIGateway(handler http.Handler) {
	lambda.Start(&defaultHandler{handler: handler})
}
