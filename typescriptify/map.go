package typescriptify

import (
	"reflect"
	"strings"
	"time"
	"fmt"
)

func (t *TypeScriptify) convertType(typeOf reflect.Type, customCode map[string]string) (string, []reflect.Type, error) {
	types := []reflect.Type{}
	if _, found := t.alreadyConverted[typeOf]; found { // Already converted
		return "", types, nil
	}

	if _, found := t.Exclude[typeOf]; found { // Exist in exclude
		return "", types, nil
	}

	if reflect.TypeOf(time.Now()) == typeOf {
		return "", types, nil
	}

	types = append(types, typeOf)

	t.alreadyConverted[typeOf] = true

	entityName := fmt.Sprintf("%s%s%s", t.Prefix, t.Suffix, typeOf.Name())
	result := fmt.Sprintf("class %s {\n", entityName)

	if !t.DontExport {
		result = "export " + result
	}

	builder := typeScriptClassBuilder{
		types:  t.types,
		indent: t.Indent,
		imports: make(map[string]string),
	}

	fields := deepFields(typeOf)
	for _, field := range fields {
		// Pointer
		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}

		//Field name
		jsonTag := field.Tag.Get("json")
		jsonFieldName := ""

		if len(jsonTag) > 0 {
			jsonTagParts := strings.Split(jsonTag, ",")
			if len(jsonTagParts) > 0 {
				jsonFieldName = strings.Trim(jsonTagParts[0], t.Indent)
			}
		}
		if len(jsonFieldName) <= 0 || jsonFieldName == "-" {
			continue
		}

		var err error

		// Custom
		customTransformation := field.Tag.Get(tsTransformTag)
		if customTransformation != "" {
			err = builder.AddSimpleField(jsonFieldName, field)
			continue
		}

		//Fields
		switch field.Type.Kind() {
		case reflect.Struct:
			name := field.Type.Name()

			typeScriptChunk, _, err := t.convertType(field.Type, customCode)
			if err != nil {
				return "", types, err
			}

			if field.Type == reflect.TypeOf(time.Now()) {
				name = "Date"
			} else {
				if typeScriptChunk != "" {
					result = typeScriptChunk + "\n" + result
					types = append(types, field.Type)
				} else {
					builder.AddImport(field.Type.Name(), t.typesPathes[field.Type])
				}
			}
			builder.AddStructField(jsonFieldName, name)
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Struct { // Slice of structs:
				typeScriptChunk, _, err := t.convertType(field.Type.Elem(), customCode)
				if err != nil {
					return "", types, err
					types = append(types, field.Type.Elem())
				} else {
					builder.AddImport(field.Type.Elem().Name(), t.typesPathes[field.Type.Elem()])
				}
				result = typeScriptChunk + "\n" + result
				builder.AddArrayOfStructsField(jsonFieldName, field.Type.Elem().Name())
			} else { // Slice of simple fields:
				err = builder.AddSimpleArrayField(jsonFieldName, field.Type.Elem().Name(), field.Type.Elem().Kind())
			}
		default:
			err = builder.AddSimpleField(jsonFieldName, field)
		}

		if err != nil {
			return "", types, err
		}
	}

	// Add fields text
	result += builder.fields

	// Create from method
	if t.CreateFromMethod {
		result += fmt.Sprintf("\n%sstatic createFrom(source: any) {\n", t.Indent)
		result += fmt.Sprintf("%s%sif ('string' === typeof source) source = JSON.parse(source);\n", t.Indent, t.Indent)
		result += fmt.Sprintf("%s%sconst result = new %s();\n", t.Indent, t.Indent, entityName)
		result += builder.createFromMethodBody
		result += fmt.Sprintf("%s%sreturn result;\n", t.Indent, t.Indent)
		result += fmt.Sprintf("%s}\n\n", t.Indent)
	}

	// Restore custom code
	if customCode != nil {
		code := customCode[entityName]
		result += t.Indent + "//[" + entityName + ":]\n" + code + "\n\n" + t.Indent + "//[end]\n"
	}

	result += "}"

	// Add imports
	result = builder.GetImportsString() + "\n" + result

	return result, types, nil
}


func deepFields(typeOf reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	if typeOf.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)

		kind := f.Type.Kind()
		if f.Anonymous && kind == reflect.Struct {
			//fmt.Println(v.Interface())
			fields = append(fields, deepFields(f.Type)...)
		} else {
			fields = append(fields, f)
		}
	}

	return fields
}