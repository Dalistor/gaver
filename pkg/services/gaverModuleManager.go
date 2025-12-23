package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Dalistor/gaver/pkg/types"
)

func ReadGaverModuleFile() (*types.GaverModuleFile, error) {
	moduleFile, err := os.ReadFile("gaverModule.json")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo gaverModule.json: %w", err)
	}

	var gaverModuleFile types.GaverModuleFile
	err = json.Unmarshal(moduleFile, &gaverModuleFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar arquivo gaverModule.json: %w", err)
	}

	return &gaverModuleFile, nil
}

func SetGaverModuleFile(gaverModuleFile *types.GaverModuleFile) error {
	moduleFile, err := json.Marshal(gaverModuleFile)
	if err != nil {
		return fmt.Errorf("erro ao serializar arquivo gaverModule.json: %w", err)
	}

	err = os.WriteFile("gaverModule.json", moduleFile, 0644)
	if err != nil {
		return fmt.Errorf("erro ao escrever arquivo gaverModule.json: %w", err)
	}

	return nil
}
