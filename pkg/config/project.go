package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProjectType representa o tipo de projeto
type ProjectType string

const (
	ProjectTypeServer  ProjectType = "server"
	ProjectTypeAndroid ProjectType = "android"
	ProjectTypeDesktop ProjectType = "desktop"
	ProjectTypeWeb     ProjectType = "web"
)

// ProjectConfig representa a configuração do projeto Gaver
type ProjectConfig struct {
	ProjectName string      `json:"projectName"`
	Type        ProjectType `json:"type"`
	Database    string      `json:"database,omitempty"`
	ServerPort  string      `json:"serverPort,omitempty"`
	FrontendDir string      `json:"frontendDir,omitempty"`
}

// ReadProjectConfig lê o arquivo GaverProject.json do diretório atual
func ReadProjectConfig() (*ProjectConfig, error) {
	configPath := "GaverProject.json"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler GaverProject.json: %w", err)
	}

	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao parsear GaverProject.json: %w", err)
	}

	return &config, nil
}

// WriteProjectConfig escreve o arquivo GaverProject.json
func WriteProjectConfig(config *ProjectConfig, projectPath string) error {
	configPath := filepath.Join(projectPath, "GaverProject.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever GaverProject.json: %w", err)
	}

	return nil
}

// IsValidProjectType verifica se o tipo de projeto é válido
func IsValidProjectType(projectType string) bool {
	switch ProjectType(projectType) {
	case ProjectTypeServer, ProjectTypeAndroid, ProjectTypeDesktop, ProjectTypeWeb:
		return true
	default:
		return false
	}
}
