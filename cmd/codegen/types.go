package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	xg "github.com/kerdokurs/xata-go"
	"github.com/stoewer/go-strcase"
)

type fieldType string

const (
	stringField fieldType = "string"
	intField    fieldType = "*int"
	floatField  fieldType = "*float64"
	boolField   fieldType = "*bool"
	multiField  fieldType = "[]string"
	dateField   fieldType = "time.Time"
	linkField   fieldType = "$link"
)

type modelField struct {
	Name        string
	Title       string
	Type        fieldType
	OmitIfEmpty string
}

type model struct {
	Name      string
	FirstChar string
	Fields    []modelField
}

const modelTmpl = `package {{ .PackageName }}

import (
    {{ range .Imports }}
    "{{ . }}"
    {{ end }}
)

type {{ .Model.Name }} struct {
		` + "Id string `json:\"id,omitempty\"`" + `
 		{{ range .Model.Fields -}}
	` +
	"        {{ .Title }} {{ .Type }} `json:\"{{ .Name }}{{ .OmitIfEmpty }}\"`" + `
		{{ end }}
}

func ({{ .Model.FirstChar }} {{ .Model.Name }}) ID() string {
		return {{ .Model.FirstChar }}.Id
}
`

// Which models to generate from the database
// If empty, all models will be generated
var modelsToGenerate = make(map[string]bool)

func generateTypes(schema *xg.Schema) {
	for _, table := range schema.Tables {
		if len(modelsToGenerate) > 0 {
			if _, ok := modelsToGenerate[table.Name]; !ok {
				continue
			}
		}

		fileName := fmt.Sprintf("%s/%s.go", *outPath, strcase.SnakeCase(table.Name))
		f, err := os.Create(fileName)
		errExit(err)

		modelTmpl := template.Must(template.New("model").Parse(modelTmpl))
		imports := make([]string, 0)
		toModel := tableToModel(&table, &imports)

		err = modelTmpl.Execute(f, struct {
			PackageName string
			Model       model
			Imports     []string
		}{
			PackageName: *packageName,
			Model:       toModel,
			Imports:     imports,
		})

		f.Close()
		errExit(err)
	}
}

func tableToModel(t *xg.SchemaTable, imports *[]string) model {
	fields, err := xg.Map(t.Columns, func(column xg.TableColumn) (modelField, error) {
		name := column.Name
		title := strcase.UpperCamelCase(name)

		typ := toFieldType(column.Type)
		if typ == dateField {
			*imports = append(*imports, "time")
		} else if typ == linkField {
			typ = fieldType("*" + column.Link.Table)
		}

		omitIfEmpty := ""
		if !column.NotNull {
			omitIfEmpty = ",omitempty"
		}

		return modelField{
			Title:       title,
			Name:        name,
			Type:        typ,
			OmitIfEmpty: omitIfEmpty,
		}, nil
	})
	if err != nil {
		panic(err)
	}

	return model{
		Name:      t.Name,
		FirstChar: strings.ToLower(t.Name[:1]),
		Fields:    fields,
	}
}

func toFieldType(columnType string) fieldType {
	var typ fieldType
	switch columnType {
	case "int":
		typ = intField
	case "float":
		typ = floatField
	case "bool":
		typ = boolField
	case "multiple":
		typ = multiField
	case "datetime":
		typ = dateField
	case "link":
		typ = linkField
	case "text":
		fallthrough
	case "string":
		fallthrough
	case "email":
		fallthrough
	default:
		typ = stringField
	}
	return typ
}

func parseModelsToGenerate() {
	if *onlyModels == "" {
		return
	}

	for _, model := range strings.Split(*onlyModels, ",") {
		modelsToGenerate[model] = true
	}
}
