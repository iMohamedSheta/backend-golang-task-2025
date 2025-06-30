package load

import (
	"taskgo/internal/config"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"
)

func InitConfig(c *ioc.Container) {
	err := ioc.Singleton(c, func(c *ioc.Container) (*config.Config, error) {
		cfg := config.New()                // Global config
		config.ApplyRegisteredLoaders(cfg) // Apply config module registered loaders
		return cfg, nil
	})

	if err != nil {
		utils.PrintErr("Failed to load config module as singleton in the ioc container: " + err.Error())
	}
}
