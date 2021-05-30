//+build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/mostafasolati/leviathan/contracts"
	"github.com/mostafasolati/leviathan/hello"
)

func InitProject(name2 string) contracts.IGreeter {
	wire.Build(hello.NewGreeter, hello.NewMostafa)
	return nil
}
