package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	coreinterfaces "github.com/silviomfa/go-cloud-core/pkg/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
)

// DynamoDBProvider implementa a interface coreinterfaces.StorageProvider para DynamoDB
type DynamoDBProvider struct {
	client   *dynamodb.Client
	provider *provider.Provider
}

// NewDynamoDBProvider cria um novo provedor de armazenamento DynamoDB
func NewDynamoDBProvider(cloudProvider coreinterfaces.CloudProvider) (coreinterfaces.StorageProvider, error) {
	log.Printf("Inicializando provedor DynamoDB com SDK v2")
	
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
	
	log.Printf("Provedor DynamoDB inicializado com sucesso")

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
	log.Printf("DynamoDB GetItem: tabela=%s, chave=%+v", tableName, key)
	
	// Converter chave para formato DynamoDB
	keyAttr, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Printf("Erro ao converter chave para atributos do DynamoDB: %v", err)
		return fmt.Errorf("erro ao converter chave para atributos do DynamoDB: %w", err)
	}
	
	log.Printf("Chave convertida: %+v", keyAttr)

	// Buscar item no DynamoDB
	response, err := p.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       keyAttr,
	})
	if err != nil {
		log.Printf("Erro ao buscar item no DynamoDB: %v", err)
		return fmt.Errorf("erro ao buscar item no DynamoDB: %w", err)
	}

	// Verificar se o item foi encontrado
	if response.Item == nil {
		log.Printf("Item não encontrado na tabela %s", tableName)
		return nil
	}
	
	log.Printf("Item encontrado: %+v", response.Item)

	// Converter item para o tipo de resultado
	if err := attributevalue.UnmarshalMap(response.Item, result); err != nil {
		log.Printf("Erro ao converter item do DynamoDB: %v", err)
		return fmt.Errorf("erro ao converter item do DynamoDB: %w", err)
	}
	
	log.Printf("Item convertido com sucesso")

	return nil
}

// PutItem insere um item no DynamoDB
func (p *DynamoDBProvider) PutItem(ctx context.Context, tableName string, item interface{}) error {
	log.Printf("DynamoDB PutItem: tabela=%s", tableName)
	
	// Converter item para formato DynamoDB
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Printf("Erro ao converter item para atributos do DynamoDB: %v", err)
		return fmt.Errorf("erro ao converter item para atributos do DynamoDB: %w", err)
	}
	
	log.Printf("Item convertido: %+v", av)

	// Inserir item no DynamoDB
	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		log.Printf("Erro ao inserir item no DynamoDB: %v", err)
		return fmt.Errorf("erro ao inserir item no DynamoDB: %w", err)
	}
	
	log.Printf("Item inserido com sucesso na tabela %s", tableName)

	return nil
}

// DeleteItem remove um item do DynamoDB
func (p *DynamoDBProvider) DeleteItem(ctx context.Context, tableName string, key map[string]interface{}) error {
	log.Printf("DynamoDB DeleteItem: tabela=%s, chave=%+v", tableName, key)
	
	// Converter chave para formato DynamoDB
	keyAttr, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Printf("Erro ao converter chave para atributos do DynamoDB: %v", err)
		return fmt.Errorf("erro ao converter chave para atributos do DynamoDB: %w", err)
	}
	
	log.Printf("Chave convertida: %+v", keyAttr)

	// Remover item do DynamoDB
	_, err = p.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       keyAttr,
	})
	if err != nil {
		log.Printf("Erro ao remover item do DynamoDB: %v", err)
		return fmt.Errorf("erro ao remover item do DynamoDB: %w", err)
	}
	
	log.Printf("Item removido com sucesso da tabela %s", tableName)

	return nil
}

// Query consulta itens no DynamoDB
func (p *DynamoDBProvider) Query(ctx context.Context, tableName string, keyCondition string, values map[string]interface{}) ([]map[string]interface{}, error) {
	log.Printf("DynamoDB Query: tabela=%s, condição=%s", tableName, keyCondition)
	
	// Se não houver condição, usar Scan em vez de Query
	if keyCondition == "" {
		return p.scan(ctx, tableName)
	}
	
	// Converter valores para formato DynamoDB
	expressionValues := make(map[string]types.AttributeValue)
	for k, v := range values {
		av, err := attributevalue.Marshal(v)
		if err != nil {
			log.Printf("Erro ao converter valor para atributo do DynamoDB: %v", err)
			return nil, fmt.Errorf("erro ao converter valor para atributo do DynamoDB: %w", err)
		}
		expressionValues[":"+k] = av
	}
	
	log.Printf("Valores convertidos: %+v", expressionValues)

	// Executar consulta no DynamoDB
	response, err := p.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expressionValues,
	})
	if err != nil {
		log.Printf("Erro ao consultar itens no DynamoDB: %v", err)
		return nil, fmt.Errorf("erro ao consultar itens no DynamoDB: %w", err)
	}
	
	log.Printf("Consulta executada com sucesso, %d itens encontrados", len(response.Items))

	// Converter itens para o formato de resultado
	result := make([]map[string]interface{}, len(response.Items))
	for i, item := range response.Items {
		m := make(map[string]interface{})
		if err := attributevalue.UnmarshalMap(item, &m); err != nil {
			log.Printf("Erro ao converter item do DynamoDB: %v", err)
			return nil, fmt.Errorf("erro ao converter item do DynamoDB: %w", err)
		}
		result[i] = m
	}
	
	log.Printf("Itens convertidos com sucesso")

	return result, nil
}

// scan executa um Scan no DynamoDB
func (p *DynamoDBProvider) scan(ctx context.Context, tableName string) ([]map[string]interface{}, error) {
	log.Printf("DynamoDB Scan: tabela=%s", tableName)

	// Executar scan no DynamoDB
	response, err := p.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Erro ao listar itens no DynamoDB: %v", err)
		return nil, fmt.Errorf("erro ao listar itens no DynamoDB: %w", err)
	}
	
	log.Printf("Scan executado com sucesso, %d itens encontrados", len(response.Items))

	// Converter itens para o formato de resultado
	result := make([]map[string]interface{}, len(response.Items))
	for i, item := range response.Items {
		m := make(map[string]interface{})
		if err := attributevalue.UnmarshalMap(item, &m); err != nil {
			log.Printf("Erro ao converter item do DynamoDB: %v", err)
			return nil, fmt.Errorf("erro ao converter item do DynamoDB: %w", err)
		}
		result[i] = m
	}
	
	log.Printf("Itens convertidos com sucesso")

	return result, nil
}