package bootstrap

import (
	"taskgo/bootstrap/load"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"

	"github.com/joho/godotenv"
)

type appBuilder struct {
	envFile   string
	actions   []func()
	container *ioc.Container
}

// Start a new builder
func NewAppBuilder(envFile string) *appBuilder {
	return &appBuilder{
		envFile:   envFile,
		container: ioc.NewContainer(),
	}
}

// Add a custom step
func (b *appBuilder) Then(fn func()) *appBuilder {
	b.actions = append(b.actions, fn)
	return b
}

// Load environment + config
func (b *appBuilder) LoadConfig() *appBuilder {
	b.mustLoadEnvFile()
	ioc.SetAppContainer(b.container)
	load.InitConfig(b.container)
	b.runActions()
	return b
}

func (b *appBuilder) LoadLogger() *appBuilder {
	load.InitLogger(b.container)
	b.runActions()
	return b
}

func (b *appBuilder) LoadDatabase() *appBuilder {
	load.InitDatabase(b.container)
	b.runActions()
	return b
}

func (b *appBuilder) LoadValidator() *appBuilder {
	load.InitValidator(b.container, registeredRules)
	b.runActions()
	return b
}

func (b *appBuilder) LoadRedisCache() *appBuilder {
	load.InitRedisCache(b.container)
	b.runActions()
	return b
}

func (b *appBuilder) LoadRedisQueue() *appBuilder {
	load.InitRedisQueue(b.container)
	b.runActions()
	return b
}

func (b *appBuilder) LoadNotify() *appBuilder {
	load.InitNotify(b.container, registerNotifyChannelsHandlers())
	b.runActions()
	return b
}

func (b *appBuilder) LoadWebsocketServer() *appBuilder {
	load.InitWebsocketServer(b.container)
	b.runActions()
	registerWebSocketsChannels(b.container)
	return b
}

func (b *appBuilder) runActions() {
	for _, fn := range b.actions {
		fn()
	}
	b.actions = nil // clear after run
}

func (b *appBuilder) mustLoadEnvFile() {
	if err := godotenv.Load(b.envFile); err != nil {
		utils.PrintErr("Error loading environment file:", err)
	}
	utils.PrintSuccess("Loaded environment file:" + b.envFile)
}

func (b *appBuilder) Boot() *appBuilder {
	registerServiceProviders(b.container)
	b.container.Bootstrap()
	b.runActions()
	return b
}
