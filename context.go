package httplam

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

type contextKey string

const (
	lambdaRequestContextKey contextKey = "lambda_request"
)

func setLambdaRequest(ctx context.Context, req *events.APIGatewayV2HTTPRequest) context.Context {
	return context.WithValue(ctx, lambdaRequestContextKey, req)
}

func GetLambdaRequest(ctx context.Context) (req *events.APIGatewayV2HTTPRequest) {
	req, _ = ctx.Value(lambdaRequestContextKey).(*events.APIGatewayV2HTTPRequest)
	return
}