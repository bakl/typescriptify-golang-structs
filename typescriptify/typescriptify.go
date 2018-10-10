package typescriptify

import (
	"reflect"
)

const (
	tsTransformTag = "ts_transform"
	tsType         = "ts_type"
)

type TypeScriptify struct {
	Prefix           string
	Suffix           string
	Indent           string
	CreateFromMethod bool
	BackupDir        string // If empty no backup
	DontExport       bool
	Exclude          map[reflect.Type]bool

	golangTypes []reflect.Type
	pathes      []string
	types       map[reflect.Kind]string

	// throwaway, used when converting
	alreadyConverted map[reflect.Type]bool
	typesPathes map[reflect.Type]string
}

func New() *TypeScriptify {
	result := new(TypeScriptify)
	result.Indent = "\t"
	result.BackupDir = "."

	types := make(map[reflect.Kind]string)

	types[reflect.Bool] = "boolean"
	types[reflect.Interface] = "any"

	types[reflect.Int] = "number"
	types[reflect.Int8] = "number"
	types[reflect.Int16] = "number"
	types[reflect.Int32] = "number"
	types[reflect.Int64] = "number"
	types[reflect.Uint] = "number"
	types[reflect.Uint8] = "number"
	types[reflect.Uint16] = "number"
	types[reflect.Uint32] = "number"
	types[reflect.Uint64] = "number"
	types[reflect.Float32] = "number"
	types[reflect.Float64] = "number"

	types[reflect.String] = "string"

	result.types = types

	result.Indent = "    "
	result.CreateFromMethod = true

	return result
}

func (t *TypeScriptify) Add(obj interface{}, path string) {
	t.AddType(reflect.TypeOf(obj), path)
}

func (t *TypeScriptify) AddType(typeOf reflect.Type, path string) {
	t.golangTypes = append(t.golangTypes, typeOf)
	t.pathes = append(t.pathes, path)
}