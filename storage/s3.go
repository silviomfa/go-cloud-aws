package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"crypto/sha256"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	coreinterfaces "github.com/silviomfa/go-cloud-core/pkg/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
)

// S3Provider implementa a interface coreinterfaces.StorageProvider para S3
type S3Provider struct {
	client   *s3.Client
	provider *provider.Provider
}

// NewS3Provider cria um novo provedor de armazenamento S3
func NewS3Provider(cloudProvider coreinterfaces.CloudProvider) (*S3Provider, error) {
	awsProvider, ok := cloudProvider.(*provider.Provider)
	if !ok {
		return nil, fmt.Errorf("provedor não é do tipo AWS")
	}

	awsConfig, ok := awsProvider.GetConfig().(aws.Config)
	if !ok {
		return nil, fmt.Errorf("configuração não é do tipo AWS")
	}

	client := s3.NewFromConfig(awsConfig)

	return &S3Provider{
		client:   client,
		provider: awsProvider,
	}, nil
}

// GetName retorna o nome do provedor
func (p *S3Provider) GetName() string {
	return "AWS-S3"
}

// GetItem recupera um objeto do S3
// Para S3, o key deve conter uma chave "Key" com o caminho do objeto
func (p *S3Provider) GetItem(ctx context.Context, bucketName string, key map[string]interface{}, result interface{}) error {
	// Extrair a chave do objeto
	objectKey, ok := key["Key"]
	if !ok {
		return fmt.Errorf("chave 'Key' não encontrada no mapa de chaves")
	}
	
	keyStr, ok := objectKey.(string)
	if !ok {
		return fmt.Errorf("chave 'Key' não é uma string")
	}
	
	// Obter objeto do S3
	output, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyStr),
	})
	if err != nil {
		return fmt.Errorf("erro ao obter objeto do S3: %w", err)
	}
	defer output.Body.Close()
	
	// Ler o conteúdo do objeto
	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler conteúdo do objeto: %w", err)
	}
	
	// Se result for um ponteiro para []byte, atribuir diretamente
	if byteSlice, ok := result.(*[]byte); ok {
		*byteSlice = data
		return nil
	}
	
	// Caso contrário, tentar deserializar como JSON
	return json.Unmarshal(data, result)
}

// PutItem insere um objeto no S3
// Para S3, o item deve ser um mapa com "Key" e "Content"
func (p *S3Provider) PutItem(ctx context.Context, bucketName string, item interface{}) error {
	var key string
	var content []byte
	
	// Verificar o tipo do item
	switch v := item.(type) {
	case map[string]interface{}:
		// Extrair a chave e o conteúdo
		keyObj, ok := v["Key"]
		if !ok {
			return fmt.Errorf("chave 'Key' não encontrada no item")
		}
		key, ok = keyObj.(string)
		if !ok {
			return fmt.Errorf("chave 'Key' não é uma string")
		}
		
		contentObj, ok := v["Content"]
		if !ok {
			return fmt.Errorf("chave 'Content' não encontrada no item")
		}
		
		// Converter conteúdo para []byte
		switch c := contentObj.(type) {
		case []byte:
			content = c
		case string:
			content = []byte(c)
		default:
			var err error
			content, err = json.Marshal(contentObj)
			if err != nil {
				return fmt.Errorf("erro ao serializar conteúdo: %w", err)
			}
		}
	default:
		// Se não for um mapa, serializar o item inteiro como JSON
		var err error
		content, err = json.Marshal(item)
		if err != nil {
			return fmt.Errorf("erro ao serializar item: %w", err)
		}
		
		// Usar um nome de arquivo baseado no hash do conteúdo
		h := sha256.New()
		h.Write(content)
		key = fmt.Sprintf("%x.json", h.Sum(nil))
	}
	
	// Inserir objeto no S3
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	})
	return err
}

// DeleteItem remove um objeto do S3
func (p *S3Provider) DeleteItem(ctx context.Context, bucketName string, key map[string]interface{}) error {
	// Extrair a chave do objeto
	objectKey, ok := key["Key"]
	if !ok {
		return fmt.Errorf("chave 'Key' não encontrada no mapa de chaves")
	}
	
	keyStr, ok := objectKey.(string)
	if !ok {
		return fmt.Errorf("chave 'Key' não é uma string")
	}
	
	// Remover objeto do S3
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyStr),
	})
	return err
}

// Query lista objetos no S3 com um prefixo
func (p *S3Provider) Query(ctx context.Context, bucketName string, keyCondition string, values map[string]interface{}) ([]map[string]interface{}, error) {
	// Para S3, keyCondition é interpretado como um prefixo
	prefix := keyCondition
	
	// Listar objetos no S3
	output, err := p.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar objetos do S3: %w", err)
	}
	
	// Converter objetos para o formato de resultado
	result := make([]map[string]interface{}, len(output.Contents))
	for i, obj := range output.Contents {
		result[i] = map[string]interface{}{
			"Key":          *obj.Key,
			"Size":         obj.Size,
			"LastModified": obj.LastModified,
			"ETag":         *obj.ETag,
		}
	}
	
	return result, nil
}
