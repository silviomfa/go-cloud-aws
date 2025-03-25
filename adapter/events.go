package adapter

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/silviomfa/go-cloud-aws/interfaces"
)

// ConvertToGenericEvent converte um evento AWS para o formato genérico
func ConvertToGenericEvent(ctx context.Context, rawEvent json.RawMessage) (interfaces.Event, error) {
	// Evento padrão
	event := interfaces.Event{}

	// Tentar identificar o tipo de evento
	var apiGatewayEvent events.APIGatewayProxyRequest
	if err := json.Unmarshal(rawEvent, &apiGatewayEvent); err == nil && apiGatewayEvent.RequestContext.RequestID != "" {
		// É um evento de API Gateway
		return convertAPIGatewayEvent(apiGatewayEvent), nil
	}

	var sqsEvent events.SQSEvent
	if err := json.Unmarshal(rawEvent, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 {
		// É um evento de SQS
		return convertSQSEvent(sqsEvent), nil
	}

	// Se não conseguir identificar, retorna o evento bruto
	event.Type = "unknown"
	event.Data = rawEvent

	return event, nil
}

// Converter evento de API Gateway
func convertAPIGatewayEvent(apiEvent events.APIGatewayProxyRequest) interfaces.Event {
	event := interfaces.Event{
		ID:        apiEvent.RequestContext.RequestID,
		Source:    "api",
		Type:      "http.request",
		RequestID: apiEvent.RequestContext.RequestID,
		Metadata: map[string]string{
			"method":     apiEvent.HTTPMethod,
			"path":       apiEvent.Path,
			"sourceIP":   apiEvent.RequestContext.Identity.SourceIP,
			"userAgent":  apiEvent.RequestContext.Identity.UserAgent,
		},
	}

	// Converter corpo para bytes
	event.Data = []byte(apiEvent.Body)

	return event
}

// Converter evento de SQS
func convertSQSEvent(sqsEvent events.SQSEvent) interfaces.Event {
	event := interfaces.Event{
		Source: "sqs",
		Type:   "message.received",
	}

	if len(sqsEvent.Records) > 0 {
		event.ID = sqsEvent.Records[0].MessageId
		event.Data = []byte(sqsEvent.Records[0].Body)
		event.Metadata = map[string]string{
			"queueUrl": sqsEvent.Records[0].EventSourceARN,
		}
	}

	return event
}

// ConvertToAWSResponse converte uma resposta genérica para o formato AWS
func ConvertToAWSResponse(response interface{}) interface{} {
	// Se já for uma resposta de API Gateway, retornar diretamente
	if apiResp, ok := response.(events.APIGatewayProxyResponse); ok {
		return apiResp
	}

	// Tentar converter para JSON
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		// Em caso de erro, retornar erro 500
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Internal Server Error"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	}

	// Retornar como resposta de API Gateway
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}