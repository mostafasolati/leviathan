package leviathan

import (
	"log"

	"github.com/mostafasolati/leviathan/config"
	"github.com/mostafasolati/leviathan/contracts"
	"github.com/mostafasolati/leviathan/logger"
	server "github.com/mostafasolati/leviathan/server"
)

type leviathan struct {
	config          contracts.IConfigService
	logger          contracts.ILogger
	serverContainer contracts.IServerContainer
}

func NewLeviathan(filename, configPath string) contracts.ILeviathan {
	configService := config.NewConfigService(configPath, filename)
	logger, err := logger.NewLogger(configService)
	if err != nil {
		log.Fatal(err)
	}

	return &leviathan{
		config: configService,
		logger: logger,
	}
}

func (s *leviathan) Config() contracts.IConfigService {
	return s.config
}

func (s *leviathan) ServerContainer() contracts.IServerContainer {
	if s.serverContainer == nil {
		s.serverContainer = server.NewEchoServerContainer(s.config, s.logger)
	}
	return s.serverContainer
}
