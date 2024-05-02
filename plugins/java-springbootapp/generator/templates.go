package generator

import _ "embed"

//go:embed templates/pom.xml.tmpl
var PomXmlTemplate string

//go:embed templates/src/main/java/Application.java.tmpl
var MainJavaTemplate string

//go:embed templates/src/main/java/HelloController.java.tmpl
var HelloControllerJavaTemplate string
