# Xata Go Codegen

This is a code generator for Go. It is used to generate code from Xata database.

## Usage

Add a file in the root of your project (usually named xata.go):

```go
//go:generate go run github.com/kerdokurs/xata-go/cmd/codegen -out database -package database
package main
```

Then run `go generate` in the root of your project.

## Options

```
Usage of codegen:
  -out string
        Output directory (default "xata")
  -package string
        Package name (default "xata")
```
