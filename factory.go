package aws

import (
	"github.com/silviomfa/go-cloud-aws/interfaces"
	"github.com/silviomfa/go-cloud-aws/provider"
	"github.com/silviomfa/go-cloud-aws/runtime"
	"github.com/silviomfa/go-cloud-aws/storage"
)

// NewProvider cria um novo provedor AWS
func NewProvider() (*provider.Provider, error) {
	return provider.NewProvider()
}

// NewDynamoDBProvider cria um novo provedor de armazenamento DynamoDB
func NewDynamoDBProvider(cloudProvider interfaces.CloudProvider) (interfaces.StorageProvider, error) {
	return storage.NewDynamoDBProvider(cloudProvider)
}

// NewRuntime cria um novo runtime AWS Lambda
func NewRuntime() interfaces.Runtime {
	return runtime.NewLambdaRuntime()
}