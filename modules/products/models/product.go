package models

import (
	"time"
)

// Product representa o model Product
// 
// COMO USAR ANNOTATIONS gaverModel:
// 
// Controle de Acesso:
//   writable:post,put,patch  - Define quais métodos HTTP podem escrever este campo
//   readable                 - Campo pode ser lido em GET
//   ignore                   - Campo completamente ignorado na API
//   ignore:write             - Nunca pode ser escrito
//   ignore:read              - Nunca pode ser lido
// 
// Validações:
//   required                 - Campo obrigatório
//   unique                   - Valor único no banco
//   email                    - Valida formato de email
//   url                      - Valida formato de URL
//   min:N                    - Valor mínimo (números)
//   max:N                    - Valor máximo (números)
//   minLength:N              - Tamanho mínimo (strings)
//   maxLength:N              - Tamanho máximo (strings)
//   pattern:regex            - Validação por regex
//   enum:val1,val2,val3      - Valores permitidos
// 
// Relacionamentos:
//   relation:hasOne          - Relação 1:1
//   relation:hasMany         - Relação 1:N
//   relation:belongsTo       - Relação N:1
//   relation:manyToMany      - Relação N:N
//   foreignKey:field_name    - Chave estrangeira
//   through:table_name       - Tabela intermediária (M2M)
// 
// EXEMPLOS:
// 
// Campo simples:
//   // gaverModel: writable:post,put; readable; required
//   Name string `json:"name"`
// 
// Campo com validação:
//   // gaverModel: writable:post; readable; required; email; unique
//   Email string `json:"email" gorm:"uniqueIndex"`
// 
// Campo numérico com range:
//   // gaverModel: writable:post,put,patch; readable; min:0; max:120
//   Age int `json:"age"`
// 
// Campo apenas leitura:
//   // gaverModel: ignore:write; readable
//   ViewCount int `json:"view_count" gorm:"default:0"`
// 
// Relacionamento:
//   // gaverModel: relation:belongsTo; foreignKey:company_id
//   CompanyID uint     `json:"company_id"`
//   Company   Company  `json:"company" gorm:"foreignKey:CompanyID"`
// 
type Product struct {
	// gaverModel: primaryKey; autoIncrement
	ID uint `json:"id" gorm:"primaryKey"`

	// TODO: Adicione seus campos aqui
	// Exemplo:
	// 
	// // gaverModel: writable:post,put,patch; readable; required; minLength:3; maxLength:100
	// Name string `json:"name" gorm:"type:varchar(100);not null"`
	// 
	// // gaverModel: writable:post; readable; required; unique; email
	// Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	// 
	// // gaverModel: writable:post,put,patch; readable; min:0
	// Age int `json:"age" gorm:"type:int"`
	// 
	// // gaverModel: writable:post,put,patch; readable
	// Status string `json:"status" gorm:"type:varchar(50);default:'active'"`
	// 
	// // gaverModel: ignore
	// Password string `json:"-" gorm:"type:varchar(255);not null"`

	// gaverModel: ignore:write; readable
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// gaverModel: ignore:write; readable
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Métodos helper personalizados
// Adicione aqui métodos para o seu model

// Validate executa validações customizadas antes de salvar
// func (m *Product) Validate() error {
// 	// Adicione suas validações aqui
// 	// Exemplo:
// 	// if m.Age < 18 {
// 	// 	return fmt.Errorf("idade deve ser maior que 18")
// 	// }
// 	return nil
// }

// BeforeSave é chamado antes de salvar no banco
// func (m *Product) BeforeSave() error {
// 	// Lógica antes de salvar
// 	// Exemplo: Hash de senha, normalização de dados, etc
// 	return nil
// }

// AfterFind é chamado após buscar do banco
// func (m *Product) AfterFind() error {
// 	// Lógica após buscar
// 	return nil
// }

// TableName especifica o nome da tabela (opcional)
// func (Product) TableName() string {
// 	return "products"
// }

