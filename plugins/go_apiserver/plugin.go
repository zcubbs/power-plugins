package main

import (
	_ "embed"
	"github.com/zcubbs/blueprint"
)

// Plugin is the exported plugin blueprint.
var Plugin blueprint.Generator = &Generator{}
