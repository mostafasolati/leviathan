package main

import "github.com/mostafasolati/leviathan"

func main() {
	lev := leviathan.Init()

	lev.Logger().Info("HELLO WORLD!")
	lev.Server().Run(":8080")
}
