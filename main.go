//+build wireinject

package leviathan

import (
	"github.com/google/wire"
	"github.com/mostafasolati/leviathan/config"
	"github.com/mostafasolati/leviathan/contracts"
	"github.com/mostafasolati/leviathan/logger"
	server "github.com/mostafasolati/leviathan/server"
)

func Init(filename string) contracts.ILeviathan {
	wire.Build(
		config.NewConfigService,
		logger.NewLogger,
		server.NewEchoServerContainer,
		NewLeviathan,
	)
	return &leviathan{}
}

type leviathan struct {
	config          contracts.IConfigService
	logger          contracts.ILogger
	serverContainer contracts.IServerContainer
}

func NewLeviathan(
	config contracts.IConfigService,
	logger contracts.ILogger,
	serverContainer contracts.IServerContainer,
) contracts.ILeviathan {
	return &leviathan{
		config: config,
		logger: logger,
	}
}

func (s *leviathan) Logger() contracts.ILogger {
	return s.logger
}

func (s *leviathan) Config() contracts.IConfigService {
	return s.config
}

func (s *leviathan) Server() contracts.IServerContainer {
	if s.serverContainer == nil {
		s.serverContainer = server.NewEchoServerContainer(s.config, s.logger)
	}
	return s.serverContainer
}
