package interfaces

import "context"

// Runtime define a interface para runtimes serverless
type Runtime interface {
	// Start inicia o runtime
	Start(handler interface{}) error

	// Wrap adapta um handler genérico para o formato específico do runtime
	Wrap(handler GenericHandler) interface{}
}

// Event representa um evento genérico
type Event struct {
	// ID único do evento
	ID string

	// Fonte do evento (ex: api, sqs, schedule)
	Source string

	// Tipo do evento (ex: http.request, message.received)
	Type string

	// Dados do evento em formato bruto
	Data []byte

	// Metadados adicionais do evento
	Metadata map[string]string

	// ID da requisição
	RequestID string
}

// GenericHandler define a interface para handlers genéricos
type GenericHandler interface {
	Handle(ctx context.Context, event Event) (interface{}, error)
}