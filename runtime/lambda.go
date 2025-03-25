package runtime

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/silviomfa/go-cloud-aws/adapter"
	"github.com/silviomfa/go-cloud-aws/interfaces"
)

// LambdaRuntime implementa a interface interfaces.Runtime para AWS Lambda
type LambdaRuntime struct{}

// NewLambdaRuntime cria um novo runtime AWS Lambda
func NewLambdaRuntime() *LambdaRuntime {
	return &LambdaRuntime{}
}

// Start inicia o runtime AWS Lambda
func (r *LambdaRuntime) Start(handler interface{}) error {
	lambda.Start(handler)
	return nil
}

// Wrap adapta um handler genérico para o formato AWS Lambda
func (r *LambdaRuntime) Wrap(handler interfaces.GenericHandler) interface{} {
	return func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		// Converter evento AWS para evento genérico
		genericEvent, err := adapter.ConvertToGenericEvent(ctx, event)
		if err != nil {
			return nil, err
		}

		// Chamar handler genérico
		response, err := handler.Handle(ctx, genericEvent)
		if err != nil {
			return nil, err
		}

		// Converter resposta genérica para resposta AWS
		return adapter.ConvertToAWSResponse(response), nil
	}
}