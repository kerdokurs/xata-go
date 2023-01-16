package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	xg "github.com/kerdokurs/xata-go"
)

func buildClient(count ...int) *xg.Client {
	accessToken := os.Getenv("XATA_API_KEY")
	if accessToken == "" {
		if len(count) == 0 {
			log.Println("Could not find 'XATA_API_KEY' in environment, trying to load from .env file")
			loadEnv()
			return buildClient(1)
		}
		fmt.Fprintf(os.Stderr, "`XATA_API_KEY` not set in environment\n")
		os.Exit(1)
	}

	dbUrl := os.Getenv("XATA_DB_URL")
	if dbUrl == "" {
		fmt.Fprintf(os.Stderr, "`XATA_DB_URL` not set in environment\n")
		os.Exit(1)
	}

	return xg.NewClient(accessToken, dbUrl)
}

const clientTmpl = `package {{ .PackageName }}

import (
	xg "github.com/kerdokurs/xata-go"
	"log"
	"os"
)

type DBClient struct {
	*xg.Client
	{{ range .Tables -}}
		{{ .Name }} xg.Table[{{ .Name }}]
	{{ end }}
}

func buildDBClient(accessToken string, databaseURL string) *DBClient {
	xataClient := xg.NewClient(accessToken, databaseURL)
	client := &DBClient{
		Client: xataClient,
		{{ range .Tables -}}
			{{ .Name }}:  xg.NewTableImpl[{{ .Name }}](xataClient, "{{ .Name }}"),
		{{ end }}
	}

	return client
}

func setup() {
	dbAccessToken := os.Getenv("XATA_API_KEY")
	if dbAccessToken == "" {
		log.Fatalf("'XATA_API_KEY' is not defined\n")
	}

	dbBaseUrl := os.Getenv("XATA_DB_URL")
	if dbBaseUrl == "" {
		log.Fatalf("'XATA_DB_URL' is not defined\n")
	}

	db = buildDBClient(dbAccessToken, dbBaseUrl)
}

var db *DBClient

func DB() *DBClient {
	if db == nil {
		setup()
	}

	return db
}
`

func generateClient(schema *xg.Schema) {
	fileName := fmt.Sprintf("%s/client.go", *outPath)
	f, err := os.Create(fileName)
	errExit(err)

	clientTmpl := template.Must(template.New("client").Parse(clientTmpl))

	err = clientTmpl.Execute(f, struct {
		PackageName string
		Tables      []xg.SchemaTable
	}{
		PackageName: *packageName,
		Tables: filter(schema.Tables, func(t xg.SchemaTable) bool {
			_, ok := modelsToGenerate[t.Name]
			return ok
		}),
	})

	f.Close()
	errExit(err)
}
