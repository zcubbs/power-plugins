package main

import _ "embed"

//go:embed files/go.mod.tmpl
var goModTemplate string

//go:embed files/main.go.tmpl
var mainGoTemplate string

//go:embed files/router/router.go.tmpl
var routerTemplate string

//go:embed files/handlers/home.go.tmpl
var homeHandlerTemplate string

//go:embed files/database/db.go.tmpl
var databaseTemplate string
