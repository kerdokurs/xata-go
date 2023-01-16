package main

import (
	"flag"
	"os"
)

var (
	client      = buildClient()
	outPath     = flag.String("out", "xata", "Path for generation output")
	packageName = flag.String("package", "xata", "Package name for generated output")
	onlyModels  = flag.String("only", "", "Only generate models for the given comma-separated list of case-sensitive table names")
)

func main() {
	flag.Parse()

	schema, err := client.GetSchema()
	errExit(err)

	ensureGenDir()

	parseModelsToGenerate()

	generateTypes(schema)
	generateClient(schema)
}

func ensureGenDir() {
	ok, err := exists(*outPath)
	errExit(err)

	if !ok {
		os.Mkdir(*outPath, 0777)
	}
}
