package interfaces

import "context"

// StorageProvider define a interface para operações com armazenamento
type StorageProvider interface {
	// GetName retorna o nome do provedor
	GetName() string

	// GetItem recupera um item
	GetItem(ctx context.Context, tableName string, key map[string]interface{}, result interface{}) error

	// PutItem insere um item
	PutItem(ctx context.Context, tableName string, item interface{}) error

	// DeleteItem remove um item
	DeleteItem(ctx context.Context, tableName string, key map[string]interface{}) error

	// Query consulta itens
	Query(ctx context.Context, tableName string, keyCondition string, values map[string]interface{}) ([]map[string]interface{}, error)
}