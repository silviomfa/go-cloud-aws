package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/silviomfa/go-cloud-aws/adapter"
	"github.com/silviomfa/go-cloud-aws/provider"
	coreinterfaces "github.com/silviomfa/go-cloud-core/pkg/interfaces"
)

// LambdaRuntime implementa a interface coreinterfaces.RuntimeProvider para AWS Lambda
type LambdaRuntime struct {
	provider *provider.Provider
}

// NewLambdaRuntime cria um novo runtime AWS Lambda
func NewLambdaRuntime(cloudProvider coreinterfaces.CloudProvider) (*LambdaRuntime, error) {
	awsProvider, ok := cloudProvider.(*provider.Provider)
	if !ok {
		return nil, fmt.Errorf("provedor não é do tipo AWS")
	}

	return &LambdaRuntime{
		provider: awsProvider,
	}, nil
}

// GetName retorna o nome do provedor
func (r *LambdaRuntime) GetName() string {
	return "AWS-Lambda"
}

// Start inicia o runtime AWS Lambda
func (r *LambdaRuntime) Start(handler interface{}) error {
	log.Println("Iniciando runtime AWS Lambda")
	lambda.Start(handler)
	return nil
}

// Wrap adapta um handler genérico para o formato AWS Lambda
func (r *LambdaRuntime) Wrap(handler coreinterfaces.GenericHandler) interface{} {
	log.Println("Adaptando handler genérico para AWS Lambda")
	return func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		log.Printf("Recebido evento Lambda: %s", string(event))
		
		// Converter evento AWS para evento genérico
		genericEvent, err := r.ParseEvent(ctx, event)
		if err != nil {
			log.Printf("Erro ao converter evento: %v", err)
			return nil, err
		}
		
		log.Printf("Evento convertido: ID=%s, Type=%s, Source=%s", genericEvent.ID, genericEvent.Type, genericEvent.Source)
		
		// Chamar handler genérico
		response, err := handler.Handle(ctx, *genericEvent)
		if err != nil {
			log.Printf("Erro no handler: %v", err)
			return nil, err
		}
		
		// Se a resposta já for do tipo Response, formatá-la
		if resp, ok := response.(*coreinterfaces.Response); ok {
			log.Printf("Formatando resposta: StatusCode=%d", resp.StatusCode)
			return r.FormatResponse(ctx, resp), nil
		}
		
		// Caso contrário, retornar a resposta como está
		log.Printf("Retornando resposta sem formatação")
		return response, nil
	}
}

// ParseEvent converte um evento específico do provedor para o formato genérico
func (r *LambdaRuntime) ParseEvent(ctx context.Context, rawEvent interface{}) (*coreinterfaces.Event, error) {
	// Converter para json.RawMessage se necessário
	var eventBytes json.RawMessage
	switch e := rawEvent.(type) {
	case json.RawMessage:
		eventBytes = e
	case []byte:
		eventBytes = json.RawMessage(e)
	default:
		var err error
		eventBytes, err = json.Marshal(rawEvent)
		if err != nil {
			return nil, err
		}
	}

	// Usar o adaptador para converter o evento
	event, err := adapter.ConvertToGenericEvent(ctx, eventBytes)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// FormatResponse converte uma resposta genérica para o formato específico do provedor
func (r *LambdaRuntime) FormatResponse(ctx context.Context, response *coreinterfaces.Response) interface{} {
	return adapter.ConvertToAWSResponse(response)
}

// GetEnvironment retorna variáveis de ambiente da função
func (r *LambdaRuntime) GetEnvironment(ctx context.Context) map[string]string {
	// Implementação simples que retorna todas as variáveis de ambiente
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	return env
}