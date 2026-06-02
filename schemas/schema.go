package schemas

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

type Schema struct {
	ID         string
	Type       reflect.Type
	JSONSchema *jsonschema.Schema
}
