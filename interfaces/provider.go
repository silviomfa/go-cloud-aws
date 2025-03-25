package interfaces

// CloudProvider define a interface para provedores de nuvem
type CloudProvider interface {
	// GetName retorna o nome do provedor
	GetName() string

	// GetRegion retorna a região configurada
	GetRegion() string

	// GetConfig retorna a configuração específica do provedor
	GetConfig() interface{}

	// IsLocal verifica se está em ambiente local
	IsLocal() bool
}
