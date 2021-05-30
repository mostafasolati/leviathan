package contracts

type ILeviathan interface {
	Config() IConfigService
	ServerContainer() IServerContainer
}
