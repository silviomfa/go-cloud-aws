package provider

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Provider implementa a interface coreinterfaces.CloudProvider para AWS
type Provider struct {
	config aws.Config
	region string
	local  bool
}

// NewProvider cria um novo provedor AWS
func NewProvider() (*Provider, error) {
	// Obter região da variável de ambiente ou usar padrão
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	// Verificar se estamos em ambiente local
	endpoint := os.Getenv("AWS_ENDPOINT")
	local := endpoint != ""

	// Configurar opções do AWS SDK
	var options []func(*config.LoadOptions) error

	// Se estiver em ambiente local, configurar endpoint personalizado
	if local {
		options = append(options, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               endpoint,
						HostnameImmutable: true,
					}, nil
				},
			),
		))
	}

	// Carregar configuração AWS
	cfg, err := config.LoadDefaultConfig(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	// Definir região
	cfg.Region = region

	return &Provider{
		config: cfg,
		region: region,
		local:  local,
	}, nil
}

// GetName retorna o nome do provedor
func (p *Provider) GetName() string {
	return "aws"
}

// GetRegion retorna a região configurada
func (p *Provider) GetRegion() string {
	return p.region
}

// GetConfig retorna a configuração específica do provedor
func (p *Provider) GetConfig() interface{} {
	return p.config
}

// IsLocal verifica se está em ambiente local
func (p *Provider) IsLocal() bool {
	return p.local
}