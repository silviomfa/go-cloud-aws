package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	coreinterfaces "github.com/silviomfa/go-cloud-core/pkg/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
)

// SQSProvider implementa a interface coreinterfaces.MessagingProvider para SQS
type SQSProvider struct {
	client   *sqs.Client
	provider *provider.Provider
}

// NewSQSProvider cria um novo provedor de mensageria SQS
func NewSQSProvider(cloudProvider coreinterfaces.CloudProvider) (*SQSProvider, error) {
	awsProvider, ok := cloudProvider.(*provider.Provider)
	if !ok {
		return nil, fmt.Errorf("provedor não é do tipo AWS")
	}

	awsConfig, ok := awsProvider.GetConfig().(aws.Config)
	if !ok {
		return nil, fmt.Errorf("configuração não é do tipo AWS")
	}

	client := sqs.NewFromConfig(awsConfig)

	return &SQSProvider{
		client:   client,
		provider: awsProvider,
	}, nil
}

// GetName retorna o nome do provedor
func (p *SQSProvider) GetName() string {
	return "AWS-SQS"
}

// SendMessage implementa o envio de mensagens para uma fila SQS
func (p *SQSProvider) SendMessage(ctx context.Context, queueName string, message interface{}) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueName,
		MessageBody: aws.String(string(messageBody)),
	})
	return err
}

// ReceiveMessages implementa a recepção de mensagens de uma fila SQS
func (p *SQSProvider) ReceiveMessages(ctx context.Context, queueName string, maxMessages int) ([]coreinterfaces.Message, error) {
	output, err := p.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &queueName,
		MaxNumberOfMessages: int32(maxMessages),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao receber mensagens: %w", err)
	}

	messages := make([]coreinterfaces.Message, len(output.Messages))
	for i, msg := range output.Messages {
		messages[i] = coreinterfaces.Message{
			ID:            *msg.MessageId,
			Body:          []byte(*msg.Body),
			ReceiptHandle: *msg.ReceiptHandle,
		}
	}

	return messages, nil
}

// DeleteMessage implementa a remoção de uma mensagem de uma fila SQS
func (p *SQSProvider) DeleteMessage(ctx context.Context, queueName string, receiptHandle string) error {
	_, err := p.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueName,
		ReceiptHandle: &receiptHandle,
	})
	return err
}
