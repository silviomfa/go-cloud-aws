package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/silviomfa/go-cloud-aws/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
)

// DynamoDBProvider implementa a interface interfaces.StorageProvider para DynamoDB
type DynamoDBProvider struct {
	client   *dynamodb.Client
	provider *provider.Provider
}

// NewDynamoDBProvider cria um novo provedor de armazenamento DynamoDB
func NewDynamoDBProvider(cloudProvider interfaces.CloudProvider) (*DynamoDBProvider, error) {
	// Verificar se o provedor é do tipo AWS
	awsProvider, ok := cloudProvider.(*provider.Provider)
	if !ok {
		return nil, fmt.Errorf("provedor não é do tipo AWS")
	}

	// Obter configuração AWS
	awsConfig, ok := awsProvider.GetConfig().(aws.Config)
	if !ok {
		return nil, fmt.Errorf("configuração não é do tipo AWS")
	}

	// Criar cliente DynamoDB
	client := dynamodb.NewFromConfig(awsConfig)

	return &DynamoDBProvider{
		client:   client,
		provider: awsProvider,
	}, nil
}

// GetName retorna o nome do provedor
func (p *DynamoDBProvider) GetName() string {
	return "AWS-DynamoDB"
}

// GetItem recupera um item do DynamoDB
func (p *DynamoDBProvider) GetItem(ctx context.Context, tableName string, key map[string]interface{}, result interface{}) error {
	// Converter chave para formato DynamoDB
	keyAttr, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("erro ao converter chave para atributos do DynamoDB: %w", err)
	}

	// Buscar item no DynamoDB
	response, err := p.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       keyAttr,
	})
	if err != nil {
		return fmt.Errorf("erro ao buscar item no DynamoDB: %w", err)
	}

	// Verificar se o item foi encontrado
	if response.Item == nil {
		return nil
	}

	// Converter item para o tipo de resultado
	if err := attributevalue.UnmarshalMap(response.Item, result); err != nil {
		return fmt.Errorf("erro ao converter item do DynamoDB: %w", err)
	}

	return nil
}

// PutItem insere um item no DynamoDB
func (p *DynamoDBProvider) PutItem(ctx context.Context, tableName string, item interface{}) error {
	// Converter item para formato DynamoDB
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("erro ao converter item para atributos do DynamoDB: %w", err)
	}

	// Inserir item no DynamoDB
	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("erro ao inserir item no DynamoDB: %w", err)
	}

	return nil
}

// DeleteItem remove um item do DynamoDB
func (p *DynamoDBProvider) DeleteItem(ctx context.Context, tableName string, key map[string]interface{}) error {
	// Converter chave para formato DynamoDB
	keyAttr, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("erro ao converter chave para atributos do DynamoDB: %w", err)
	}

	// Remover item do DynamoDB
	_, err = p.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key:       keyAttr,
	})
	if err != nil {
		return fmt.Errorf("erro ao remover item do DynamoDB: %w", err)
	}

	return nil
}

// Query consulta itens no DynamoDB
func (p *DynamoDBProvider) Query(ctx context.Context, tableName string, keyCondition string, values map[string]interface{}) ([]map[string]interface{}, error) {
	// Implementação para Query no DynamoDB
	// ...
	
	return nil, nil
}