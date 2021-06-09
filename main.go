//+build wireinject

package leviathan

import (
	"github.com/google/wire"
	"github.com/mostafasolati/leviathan/auth"
	"github.com/mostafasolati/leviathan/config"
	"github.com/mostafasolati/leviathan/contracts"
	"github.com/mostafasolati/leviathan/logger"
	"github.com/mostafasolati/leviathan/notification"
	server "github.com/mostafasolati/leviathan/server"
	"github.com/mostafasolati/leviathan/user"
)

func Init(filename string, apiKey notification.ApikeyType) contracts.ILeviathan {
	wire.Build(
		config.NewConfigService,
		notification.NewKavenegar,
		logger.NewLogger,
		server.NewEchoServerContainer,
		auth.NewAuthService,
		NewLeviathan,
		user.NewUserService,
	)
	return &leviathan{}
}

type leviathan struct {
	config          contracts.IConfigService
	logger          contracts.ILogger
	serverContainer contracts.IServerContainer
	user            contracts.IUserService
	auth            contracts.IAuth
}

func NewLeviathan(
	config contracts.IConfigService,
	logger contracts.ILogger,
	serverContainer contracts.IServerContainer,
	userService contracts.IUserService,
	auth contracts.IAuth,
) contracts.ILeviathan {
	return &leviathan{
		config: config,
		logger: logger,
		user:   userService,
		auth:   auth,
	}
}

func (s *leviathan) Auth() contracts.IAuth {
	return s.auth
}
func (s *leviathan) User() contracts.IUserService {
	return s.user
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
