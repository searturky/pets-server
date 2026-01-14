package providers

import (
	"os"

	"pets-server/internal/domain/pet"
	"pets-server/internal/domain/pet/interpreter"
	"pets-server/internal/pkg/config"
)

// ProvideSpeciesConfig 加载物种配置
func ProvideSpeciesConfig() (*config.SpeciesConfig, error) {
	path := "configs/species.yaml"
	if envPath := os.Getenv("SPECIES_CONFIG_PATH"); envPath != "" {
		path = envPath
	}
	return config.LoadSpecies(path)
}

// ProvideInterpreterFactory 提供解释器工厂
func ProvideInterpreterFactory() *interpreter.InterpreterFactory {
	return interpreter.NewInterpreterFactory()
}

// ProvideSpeciesRegistry 构建物种注册表
func ProvideSpeciesRegistry(
	cfg *config.SpeciesConfig,
	factory *interpreter.InterpreterFactory,
) (*pet.SpeciesRegistry, error) {
	return interpreter.BuildSpeciesRegistry(cfg, factory)
}

// ProvideFusionRegistry 构建融合注册表
func ProvideFusionRegistry(cfg *config.SpeciesConfig) *pet.SpeciesFusionRegistry {
	return interpreter.BuildFusionRegistry(cfg)
}
