package contracts

type ILeviathan interface {
	Config() IConfigService
	Server() IServerContainer
	Logger() ILogger
	Auth() IAuth
	User() IUserService
}
