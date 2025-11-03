package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Dalistor/gaver/pkg/parser"
)

// Validator valida dados baseado nas annotations do model
type Validator struct {
	metadata *parser.ModelMetadata
}

// NewValidator cria um novo validador
func NewValidator(metadata *parser.ModelMetadata) *Validator {
	return &Validator{metadata: metadata}
}

// Validate valida os dados fornecidos
func (v *Validator) Validate(data map[string]interface{}) error {
	for _, field := range v.metadata.Fields {
		value, exists := data[field.Name]

		// Verificar required
		if field.Required {
			if !exists || isNilOrEmpty(value) {
				return fmt.Errorf("campo '%s' é obrigatório", field.Name)
			}
		}

		// Se não existe e não é required, pular validações
		if !exists {
			continue
		}

		// Executar validações específicas
		if err := v.validateField(field, value); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateField(field parser.FieldMetadata, value interface{}) error {
	// Email
	if _, ok := field.Validations["email"]; ok {
		if str, ok := value.(string); ok {
			if !isValidEmail(str) {
				return fmt.Errorf("campo '%s' deve ser um email válido", field.Name)
			}
		}
	}

	// URL
	if _, ok := field.Validations["url"]; ok {
		if str, ok := value.(string); ok {
			if !isValidURL(str) {
				return fmt.Errorf("campo '%s' deve ser uma URL válida", field.Name)
			}
		}
	}

	// Pattern (regex)
	if pattern, ok := field.Validations["pattern"]; ok {
		if str, ok := value.(string); ok {
			if !matchesPattern(str, pattern) {
				return fmt.Errorf("campo '%s' não corresponde ao padrão esperado", field.Name)
			}
		}
	}

	// Enum
	if enumStr, ok := field.Validations["enum"]; ok {
		if !isInEnum(value, enumStr) {
			return fmt.Errorf("campo '%s' deve ser um dos valores: %s", field.Name, enumStr)
		}
	}

	// Min (para números)
	if minStr, ok := field.Validations["min"]; ok {
		if !validateMin(value, minStr) {
			return fmt.Errorf("campo '%s' deve ser maior ou igual a %s", field.Name, minStr)
		}
	}

	// Max (para números)
	if maxStr, ok := field.Validations["max"]; ok {
		if !validateMax(value, maxStr) {
			return fmt.Errorf("campo '%s' deve ser menor ou igual a %s", field.Name, maxStr)
		}
	}

	// MinLength (para strings)
	if minLenStr, ok := field.Validations["minLength"]; ok {
		if str, ok := value.(string); ok {
			minLen, _ := strconv.Atoi(minLenStr)
			if len(str) < minLen {
				return fmt.Errorf("campo '%s' deve ter no mínimo %d caracteres", field.Name, minLen)
			}
		}
	}

	// MaxLength (para strings)
	if maxLenStr, ok := field.Validations["maxLength"]; ok {
		if str, ok := value.(string); ok {
			maxLen, _ := strconv.Atoi(maxLenStr)
			if len(str) > maxLen {
				return fmt.Errorf("campo '%s' deve ter no máximo %d caracteres", field.Name, maxLen)
			}
		}
	}

	return nil
}

// FilterWritableFields filtra campos que podem ser escritos no método especificado
func (v *Validator) FilterWritableFields(data map[string]interface{}, method string) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range data {
		// Encontrar metadata do campo
		field := v.findField(key)
		if field == nil {
			// Campo não encontrado no model, ignorar
			continue
		}

		// Verificar se pode ser escrito
		if field.IsWritableInMethod(method) {
			filtered[key] = value
		}
	}

	return filtered
}

// FilterReadableFields filtra campos que podem ser lidos
func (v *Validator) FilterReadableFields(data map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range data {
		field := v.findField(key)
		if field == nil {
			continue
		}

		if field.IsReadable() {
			filtered[key] = value
		}
	}

	return filtered
}

func (v *Validator) findField(name string) *parser.FieldMetadata {
	// Buscar por nome exato
	for i := range v.metadata.Fields {
		if v.metadata.Fields[i].Name == name {
			return &v.metadata.Fields[i]
		}
	}

	// Buscar por JSON tag
	for i := range v.metadata.Fields {
		if v.metadata.Fields[i].JSONTag == name {
			return &v.metadata.Fields[i]
		}
	}

	return nil
}

// Funções auxiliares

func isNilOrEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	}

	return false
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidURL(url string) bool {
	re := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(url)
}

func matchesPattern(str, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(str)
}

func isInEnum(value interface{}, enumStr string) bool {
	enumValues := strings.Split(enumStr, ",")
	valueStr := fmt.Sprintf("%v", value)

	for _, enumValue := range enumValues {
		if strings.TrimSpace(enumValue) == strings.TrimSpace(valueStr) {
			return true
		}
	}

	return false
}

func validateMin(value interface{}, minStr string) bool {
	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		return true
	}

	switch v := value.(type) {
	case int:
		return float64(v) >= min
	case int64:
		return float64(v) >= min
	case float64:
		return v >= min
	case float32:
		return float64(v) >= min
	}

	return true
}

func validateMax(value interface{}, maxStr string) bool {
	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return true
	}

	switch v := value.(type) {
	case int:
		return float64(v) <= max
	case int64:
		return float64(v) <= max
	case float64:
		return v <= max
	case float32:
		return float64(v) <= max
	}

	return true
}
