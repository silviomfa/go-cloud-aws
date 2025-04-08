package aws

import (
	"log"
	
	"github.com/silviomfa/go-cloud-aws/messaging"
	"github.com/silviomfa/go-cloud-aws/provider"
	"github.com/silviomfa/go-cloud-aws/runtime"
	"github.com/silviomfa/go-cloud-aws/storage"
	"github.com/silviomfa/go-cloud-core/pkg/factory"
	"github.com/silviomfa/go-cloud-core/pkg/interfaces"
)

// Inicializar registra todos os provedores AWS
func init() {
	log.Println("Registrando provedores AWS...")
	
	// Registrar provedor de nuvem AWS
	factory.RegisterCloudProvider("aws", func(config map[string]interface{}) (interfaces.CloudProvider, error) {
		log.Println("Criando provedor de nuvem AWS")
		return provider.NewProvider()
	})
	
	// Registrar provedor de armazenamento DynamoDB
	factory.RegisterStorageProvider("aws", func(provider interfaces.CloudProvider) (interfaces.StorageProvider, error) {
		log.Println("Criando provedor de armazenamento DynamoDB")
		return storage.NewDynamoDBProvider(provider)
	})
	
	// Registrar provedor de mensageria SQS
	factory.RegisterMessagingProvider("aws", func(provider interfaces.CloudProvider) (interfaces.MessagingProvider, error) {
		log.Println("Criando provedor de mensageria SQS")
		return messaging.NewSQSProvider(provider)
	})
	
	// Registrar provedor runtime Lambda
	factory.RegisterRuntimeProvider("aws", func(provider interfaces.CloudProvider) (interfaces.RuntimeProvider, error) {
		log.Println("Criando provedor runtime Lambda")
		return runtime.NewLambdaRuntime(provider)
	})
	
	log.Println("Provedores AWS registrados com sucesso")
} 