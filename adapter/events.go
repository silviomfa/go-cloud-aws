package adapter

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	coreinterfaces "github.com/silviomfa/go-cloud-core/pkg/interfaces"
)

// ConvertToGenericEvent converte um evento AWS para o formato genérico
func ConvertToGenericEvent(ctx context.Context, eventBytes json.RawMessage) (coreinterfaces.Event, error) {
	log.Printf("Convertendo evento AWS para formato genérico: %s", string(eventBytes))
	
	// Criar evento genérico com valores padrão
	event := coreinterfaces.Event{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      eventBytes,
		Metadata:  make(map[string]interface{}),
	}
	
	// Tentar identificar o tipo de evento
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(eventBytes, &rawEvent); err != nil {
		log.Printf("Erro ao decodificar evento AWS: %v", err)
		return event, nil // Retornar evento com dados brutos
	}
	
	// Verificar se é um evento API Gateway
	if _, ok := rawEvent["httpMethod"]; ok {
		event.Type = "http.request"
		event.Source = "api"
		
		// Extrair ID da requisição se disponível
		if requestContext, ok := rawEvent["requestContext"].(map[string]interface{}); ok {
			if requestID, ok := requestContext["requestId"].(string); ok {
				event.ID = requestID
			}
		}
		
		// Adicionar pathParameters aos metadados
		if pathParams, ok := rawEvent["pathParameters"].(map[string]interface{}); ok {
			event.Metadata["pathParameters"] = pathParams
		}
		
		// Adicionar queryStringParameters aos metadados
		if queryParams, ok := rawEvent["queryStringParameters"].(map[string]interface{}); ok {
			event.Metadata["queryStringParameters"] = queryParams
		}
		
		// Adicionar body aos metadados
		if body, ok := rawEvent["body"].(string); ok {
			event.Metadata["body"] = body
		}
		
		// Adicionar httpMethod aos metadados
		if method, ok := rawEvent["httpMethod"].(string); ok {
			event.Metadata["httpMethod"] = method
		}
		
		log.Printf("Evento identificado como API Gateway: ID=%s", event.ID)
		return event, nil
	}
	
	// Verificar se é um evento SQS
	if records, ok := rawEvent["Records"].([]interface{}); ok {
		for _, record := range records {
			if recordMap, ok := record.(map[string]interface{}); ok {
				if eventSource, ok := recordMap["eventSource"].(string); ok && strings.Contains(strings.ToLower(eventSource), "sqs") {
					event.Type = "sqs.message"
					event.Source = "queue"
					
					// Extrair ID da mensagem se disponível
					if messageID, ok := recordMap["messageId"].(string); ok {
						event.ID = messageID
					}
					
					log.Printf("Evento identificado como SQS: ID=%s", event.ID)
					return event, nil
				}
			}
		}
	}
	
	// Se não conseguiu identificar o tipo, usar valores genéricos
	event.Type = "unknown"
	event.Source = "aws"
	log.Printf("Tipo de evento não identificado: %s", event.Type)
	
	return event, nil
}

// Converter evento de API Gateway
func convertAPIGatewayEvent(apiEvent events.APIGatewayProxyRequest) coreinterfaces.Event {
	event := coreinterfaces.Event{
		ID:        apiEvent.RequestContext.RequestID,
		Source:    "api",
		Type:      "http.request",
		RequestID: apiEvent.RequestContext.RequestID,
		Metadata: map[string]interface{}{
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
func convertSQSEvent(sqsEvent events.SQSEvent) coreinterfaces.Event {
	event := coreinterfaces.Event{
		Source: "sqs",
		Type:   "message.received",
	}

	if len(sqsEvent.Records) > 0 {
		event.ID = sqsEvent.Records[0].MessageId
		event.Data = []byte(sqsEvent.Records[0].Body)
		event.Metadata = map[string]interface{}{
			"queueUrl": sqsEvent.Records[0].EventSourceARN,
		}
	}

	return event
}

// ConvertToAWSResponse converte uma resposta genérica para o formato AWS
func ConvertToAWSResponse(response *coreinterfaces.Response) interface{} {
	// Converter para resposta API Gateway
	apiResponse := events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Headers:    make(map[string]string),
		Body:       string(response.Body),
	}
	
	// Copiar headers
	for k, v := range response.Headers {
		apiResponse.Headers[k] = v
	}
	
	return apiResponse
}