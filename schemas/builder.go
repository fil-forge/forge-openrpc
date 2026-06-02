package schemas

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alanshaw/go-openrpc"
	"github.com/google/jsonschema-go/jsonschema"
)

func BuildMapSchema(obj any, refSchemas []Schema) (Schema, []Schema, error) {
	props := map[string]*jsonschema.Schema{}
	refs := []Schema{}

	value := reflect.ValueOf(obj)
	typ := value.Type()

	optional := map[string]struct{}{}

	for i := 0; i < value.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if tag, ok := typ.Field(i).Tag.Lookup("cborgen"); ok {
			name := strings.Split(tag, ",")[0] // TODO: this is not always a rename
			if name != "" {
				fieldName = name
			}
		}

		fieldValue := value.Field(i)
		fieldType := fieldValue.Type()

		if fieldType.Kind() == reflect.Pointer {
			optional[fieldName] = struct{}{}
			fieldType = fieldType.Elem()
		}

		if match, ok := matchSchemaType(fieldType, refSchemas); ok {
			props[fieldName] = &jsonschema.Schema{
				Ref: "#/components/schemas/" + match.ID,
			}
			continue
		}

		switch fieldType.Kind() {
		case reflect.Int:
		case reflect.Int64:
			props[fieldName] = &jsonschema.Schema{
				Type:        "integer",
				Description: "A 64-bit integer.",
			}
			continue
		case reflect.Uint:
		case reflect.Uint64:
			min := float64(0)
			props[fieldName] = &jsonschema.Schema{
				Type:        "integer",
				Description: "An unsigned 64-bit integer.",
				Minimum:     &min,
			}
			continue
		case reflect.String:
			props[fieldName] = &jsonschema.Schema{
				Type: "string",
			}
			continue
		case reflect.Struct:
			schema, _, err := BuildMapSchema(fieldValue.Interface(), refSchemas)
			if err != nil {
				return Schema{}, nil, fmt.Errorf("building nested struct schema: %w", err)
			}
			props[fieldName] = &jsonschema.Schema{
				Ref: "#/components/schemas/" + schema.ID,
			}
			refs = append(refs, schema)
			continue
		case reflect.Slice:
			fieldElemType := fieldType.Elem()

			if match, ok := matchSchemaType(fieldElemType, refSchemas); ok {
				props[fieldName] = &jsonschema.Schema{
					Items: &jsonschema.Schema{
						Ref: "#/components/schemas/" + match.ID,
					},
				}
				continue
			}

			if fieldElemType.Kind() == reflect.Struct {
				schema, _, err := BuildMapSchema(reflect.New(fieldElemType).Elem().Interface(), refSchemas)
				if err != nil {
					return Schema{}, nil, fmt.Errorf("building nested struct schema: %w", err)
				}
				props[fieldName] = &jsonschema.Schema{
					Items: &jsonschema.Schema{
						Ref: "#/components/schemas/" + schema.ID,
					},
				}
				refs = append(refs, schema)
				continue
			}
		}
		return Schema{}, nil, fmt.Errorf("unsupported type: %T", fieldValue.Interface())
	}

	id := strings.ToLower(typ.Name()[:1]) + typ.Name()[1:]
	required := []string{}
	for name := range props {
		if _, ok := optional[name]; !ok {
			required = append(required, name)
		}
	}

	return Schema{
		ID:   id,
		Type: typ,
		JSONSchema: &jsonschema.Schema{
			Title:      typ.Name(),
			Type:       "object",
			Properties: props,
			Required:   required,
		},
	}, refs, nil
}

func BuildMethod(name string, desc string, args any, result any, refSchemas []Schema) (*openrpc.Method, []Schema, error) {
	var refs []Schema

	argsContentDescs := []*openrpc.ContentDescriptor{}
	if reflect.TypeOf(args) != nil {
		argSchema, argRefs, err := BuildMapSchema(args, refSchemas)
		if err != nil {
			return nil, nil, fmt.Errorf("building argument schema: %w", err)
		}
		refs = append(refs, argRefs...)

		for name, prop := range argSchema.JSONSchema.Properties {
			required := false
			for _, req := range argSchema.JSONSchema.Required {
				if req == name {
					required = true
					break
				}
			}
			argsContentDescs = append(argsContentDescs, &openrpc.ContentDescriptor{
				Name:     name,
				Schema:   &openrpc.JSONSchema{Schema: prop},
				Required: required,
			})
		}
	}

	var resultContentDesc *openrpc.ContentDescriptor
	if reflect.TypeOf(result) != nil {
		resultSchema, resultRefs, err := BuildMapSchema(result, refSchemas)
		if err != nil {
			return nil, nil, fmt.Errorf("building result schema: %w", err)
		}
		refs = append(refs, resultRefs...)
		resultContentDesc = &openrpc.ContentDescriptor{
			Name:   resultSchema.ID,
			Schema: &openrpc.JSONSchema{Schema: resultSchema.JSONSchema},
		}
	}

	return &openrpc.Method{
		Name:        name,
		Description: desc,
		Params:      argsContentDescs,
		Result:      resultContentDesc,
	}, refs, nil
}

func matchSchemaType(t reflect.Type, refs []Schema) (Schema, bool) {
	for _, rs := range refs {
		if rs.Type == t {
			return rs, true
		}
	}
	return Schema{}, false
}
