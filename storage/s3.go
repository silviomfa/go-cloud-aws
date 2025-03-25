package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/silviomfa/go-cloud-aws/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
)

// S3Provider implementa a interface interfaces.StorageProvider para S3
type S3Provider struct {
	client   *s3.Client
	provider *provider.Provider
}

// NewS3Provider cria um novo provedor de armazenamento S3
func NewS3Provider(cloudProvider interfaces.CloudProvider) (*S3Provider, error) {
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

// GetItem implementa a obtenção de um objeto do S3
func (p *S3Provider) GetItem(ctx context.Context, bucketName string, key string, writer io.Writer) error {
	output, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("erro ao obter objeto do S3: %w", err)
	}
	defer output.Body.Close()

	_, err = io.Copy(writer, output.Body)
	return err
}

// PutItem implementa a inserção de um objeto no S3
func (p *S3Provider) PutItem(ctx context.Context, bucketName string, key string, reader io.Reader) error {
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   reader,
	})
	return err
}

// DeleteItem implementa a remoção de um objeto do S3
func (p *S3Provider) DeleteItem(ctx context.Context, bucketName string, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	return err
}

// Query não é suportado para S3
func (p *S3Provider) Query(ctx context.Context, tableName string, keyCondition string, values map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("query não suportado para S3")
}
