package typescriptify

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
)

type typeScriptClassBuilder struct {
	types                map[reflect.Kind]string
	indent               string
	fields               string
	imports              map[string]string
	createFromMethodBody string
}

func (t *typeScriptClassBuilder) AddImport(fieldName string, typePath string) {
	t.imports[fieldName] = typePath
}

func (t *typeScriptClassBuilder) GetImportsString() string {
	res := ""

	for fieldName, typePath := range t.imports {
		fileName := path.Base(typePath)
		path := strings.TrimSuffix(fileName, path.Ext(fileName))

		if path == "" {
			continue
		}

		res += "import {"+fieldName+"} from './"+path+"';\n"
	}

	return res
}

func (t *typeScriptClassBuilder) AddStructField(fieldName, fieldType string) {
	t.fields += fmt.Sprintf("%s%s: %s;\n", t.indent, fieldName, fieldType)
	t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"] ? %s.createFrom(source[\"%s\"]) : null;\n", t.indent, t.indent, fieldName, fieldName, fieldType, fieldName)
}

func (t *typeScriptClassBuilder) AddArrayOfStructsField(fieldName, fieldType string) {
	t.fields += fmt.Sprintf("%s%s: %s[];\n", t.indent, fieldName, fieldType)
	t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"] ? source[\"%s\"].map(function(element) { return %s.createFrom(element); }) : null;\n", t.indent, t.indent, fieldName, fieldName, fieldName, fieldType)
}

func (t *typeScriptClassBuilder) AddMapOfStructsField(fieldName, keyType string, fieldType string) {
	t.fields += fmt.Sprintf("%s%s: {[key: %s]: %s};\n", t.indent, fieldName, keyType, fieldType)
}

func (t *typeScriptClassBuilder) AddMapOfSimpleField(fieldName, keyType string, valueKind reflect.Kind) error {
	if typeScriptType, ok := t.types[valueKind]; ok {
		if len(fieldName) > 0 {
			t.fields += fmt.Sprintf("%s%s: {[key: %s]: %s};\n", t.indent, fieldName, keyType, typeScriptType)
			//@TODO Create From Method?
			return nil
		}
	}
	return errors.New(fmt.Sprintf("cannot find type for %s (%s)", valueKind.String(), fieldName))
}

func (t *typeScriptClassBuilder) AddSimpleArrayField(fieldName, fieldType string, kind reflect.Kind) error {
	if typeScriptType, ok := t.types[kind]; ok {
		if len(fieldName) > 0 {
			t.fields += fmt.Sprintf("%s%s: %s[];\n", t.indent, fieldName, typeScriptType)
			t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"];\n", t.indent, t.indent, fieldName, fieldName)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("cannot find type for %s (%s/%s)", kind.String(), fieldName, fieldType))
}

func (t *typeScriptClassBuilder) AddSimpleField(fieldName string, field reflect.StructField) error {
	kind := field.Type.Kind()
	customTSType := field.Tag.Get(tsType)

	typeScriptType := t.types[kind]
	if len(customTSType) > 0 {
		typeScriptType = customTSType
	}

	customTransformation := field.Tag.Get(tsTransformTag)

	if len(typeScriptType) > 0 && len(fieldName) > 0 {
		t.fields += fmt.Sprintf("%s%s: %s;\n", t.indent, fieldName, typeScriptType)
		if customTransformation == "" {
			t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"];\n", t.indent, t.indent, fieldName, fieldName)
		} else {
			val := fmt.Sprintf(`source["%s"]`, fieldName)
			expression := strings.Replace(customTransformation, "__VALUE__", val, -1)
			t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = %s;\n", t.indent, t.indent, fieldName, expression)
		}
		return nil
	}

	return errors.New("Cannot find type for " + field.Type.String() + ", field: " + fieldName)
}

