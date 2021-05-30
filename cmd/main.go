package main

import "github.com/mostafasolati/leviathan"

func main() {
	lev := leviathan.Init("config")

	lev.Logger().Info("HELLO WORLD!")
	lev.Server().Run(":8080")
}
