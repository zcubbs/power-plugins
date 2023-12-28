package main

import (
	_ "embed"
	"github.com/hashicorp/go-plugin"
	"github.com/zcubbs/blueprint"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: blueprint.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"blueprint": &blueprint.GeneratorPlugin{Impl: &Generator{}},
		},
	})

	// Hang the main process as the plugin should be run in a separate process
	select {}
}
