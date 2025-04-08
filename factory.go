package aws

import (
	coreinterfaces "github.com/silviomfa/go-cloud-core/pkg/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
	"github.com/silviomfa/go-cloud-aws/runtime"
	"github.com/silviomfa/go-cloud-aws/storage"
	"github.com/silviomfa/go-cloud-aws/messaging"
)

// NewProvider cria um novo provedor AWS
func NewProvider() (*provider.Provider, error) {
	return provider.NewProvider()
}

// NewDynamoDBProvider cria um novo provedor de armazenamento DynamoDB
func NewDynamoDBProvider(cloudProvider coreinterfaces.CloudProvider) (coreinterfaces.StorageProvider, error) {
	return storage.NewDynamoDBProvider(cloudProvider)
}

// NewRuntime cria um novo runtime AWS Lambda
func NewRuntime(cloudProvider coreinterfaces.CloudProvider) (coreinterfaces.RuntimeProvider, error) {
	return runtime.NewLambdaRuntime(cloudProvider)
}

// NewS3Provider cria um novo provedor de armazenamento S3
func NewS3Provider(cloudProvider coreinterfaces.CloudProvider) (coreinterfaces.StorageProvider, error) {
	return storage.NewS3Provider(cloudProvider)
}

// NewSQSProvider cria um novo provedor de mensageria SQS
func NewSQSProvider(cloudProvider coreinterfaces.CloudProvider) (coreinterfaces.MessagingProvider, error) {
	return messaging.NewSQSProvider(cloudProvider)
}