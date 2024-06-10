package schema

import (
	_ "embed"

	"github.com/xeipuuv/gojsonschema"
)

var (
	//go:embed register.schema.json
	registerSF []byte
)

var (
	//Defines the JSON schema for a registration request.
	Register gojsonschema.JSONLoader = gojsonschema.NewBytesLoader(registerSF)
)
