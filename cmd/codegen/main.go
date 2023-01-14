package main

import (
	"flag"
	"os"
)

var (
	client      = buildClient()
	outPath     = flag.String("out", "xata", "Path for generation output")
	packageName = flag.String("package", "xata", "Package name for generated output")
)

func main() {
	flag.Parse()

	schema, err := client.GetSchema()
	errExit(err)

	ensureGenDir()

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
