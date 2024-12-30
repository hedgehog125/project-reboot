package main

import (
	"github.com/hedgehog125/project-reboot/subfns"
)

func main() {
	env := subfns.LoadEnvironmentVariables()
	_ = subfns.OpenDatabase(env)
	engine := subfns.ConfigureServer(env)
	subfns.RunServer(engine, env)
}
