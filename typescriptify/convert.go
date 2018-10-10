package typescriptify

import (
	"reflect"
	"strings"
	"os"
)

// Convert to string
func (t *TypeScriptify) Convert(customCode map[string]string) (map[reflect.Type]string, error) {
	t.alreadyConverted = make(map[reflect.Type]bool)
	t.typesPathes = make(map[reflect.Type]string)

	res := map[reflect.Type]string{}
	for _, typeof := range t.golangTypes {
		result := ""
		typeScriptCode, _, err := t.convertType(typeof, customCode)
		if err != nil {
			return res, err
		}
		result += "\n" + strings.Trim(typeScriptCode, " "+t.Indent+"\r\n")

		res[typeof] = result
	}

	return res, nil
}

// Convert to files
func (t TypeScriptify) ConvertToFiles() error {
	t.alreadyConverted = make(map[reflect.Type]bool)
	t.typesPathes = make(map[reflect.Type]string)

	for i, typeof := range t.golangTypes {
		fileName := t.pathes[i]

		result := ""

		customCode, err := loadCustomCode(fileName)
		if err != nil {
			return err
		}

		typeScriptCode, types, err := t.convertType(typeof, customCode)
		if err != nil {
			return err
		}
		result += "\n" + strings.Trim(typeScriptCode, " "+t.Indent+"\r\n")

		for _, typeOf := range types {
			t.typesPathes[typeOf] = fileName
		}

		if len(t.BackupDir) > 0 {
			err := t.backup(fileName)
			if err != nil {
				return err
			}
		}

		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		f.WriteString("/* Do not change, this code is generated from Golang structs */\n\n")
		f.WriteString(result)
		if err != nil {
			return err
		}

		f.Close()
	}

	return nil
}
