package interfaces

import "context"

// MessagingProvider define a interface para operações com mensageria
type MessagingProvider interface {
	// GetName retorna o nome do provedor
	GetName() string

	// SendMessage envia uma mensagem para uma fila
	SendMessage(ctx context.Context, queueName string, message interface{}) error

	// ReceiveMessages recebe mensagens de uma fila
	ReceiveMessages(ctx context.Context, queueName string, maxMessages int) ([]Message, error)

	// DeleteMessage remove uma mensagem da fila
	DeleteMessage(ctx context.Context, queueName string, receiptHandle string) error
}

// Message representa uma mensagem genérica
type Message struct {
	ID            string
	Body          []byte
	Attributes    map[string]string
	ReceiptHandle string
}