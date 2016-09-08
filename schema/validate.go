package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/xeipuuv/gojsonschema"
)

const usage = `check document with specified schema

Users can use validate in following ways.
validate --schema <schema.json> <document.json>`

func validate(context *cli.Context) error {
	nargs := context.NArg()

	schemaFile := context.String("schema")
	if schemaFile == "" {
		return fmt.Errorf("Error: schema-json file must be specified!")
	}
	schemaPath, err := filepath.Abs(schemaFile)

	if err != nil {
		return fmt.Errorf("Error: invalid schema-file path: %s", err)
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	var documentLoader gojsonschema.JSONLoader

	if nargs == 1 {
		documentArg := context.Args().Get(0)
		documentPath, err := filepath.Abs(documentArg)
		if err != nil {
			return fmt.Errorf("Error: invalid document-file path: %s\n", err)
		}
		documentLoader = gojsonschema.NewReferenceLoader("file://" + documentPath)
	} else if nargs == 0 {
		documentBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Error: read input document failed: %s\n", err)
		}
		documentString := string(documentBytes)
		documentLoader = gojsonschema.NewStringLoader(documentString)
	} else {
		return fmt.Errorf("Error: invalid arguments number\n")
	}

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "validate"
	app.Usage = usage

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "schema,s",
			Value: "",
			Usage: "specify the schema-json file",
		},
	}

	app.Action = validate

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
