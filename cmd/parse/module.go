package parse

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/cliutil/cmder"
	"github.com/tsaikd/KDGoLib/jsonex"
	"github.com/tsaikd/go-raml-parser/parser"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
	"gopkg.in/urfave/cli.v2"
)

// Module info
var Module = cmder.NewModule("parse").
	SetUsage("Parse RAML file and show API in json format").
	AddFlag(
		&cli.StringFlag{
			Name:        "f",
			Aliases:     []string{"ramlfile"},
			Value:       "api.raml",
			Usage:       "Source RAML file",
			Destination: &ramlFile,
		},
		&cli.BoolFlag{
			Name:        "checkRAMLVersion",
			Usage:       "Check RAML Version",
			Destination: &checkRAMLVersion,
		},
		&cli.BoolFlag{
			Name:        "ignoreUnusedAnnotation",
			Usage:       "Ignore unused annotations",
			Destination: &ignoreUnusedAnnotation,
		},
		&cli.BoolFlag{
			Name:        "ignoreUnusedTrait",
			Usage:       "Ignore unused traits",
			Destination: &ignoreUnusedTrait,
		},
		&cli.BoolFlag{
			Name:        "allowIntegerToBeNumber",
			Usage:       "Allow integer type to be number type when checking",
			Destination: &allowIntegerToBeNumber,
		},
		&cli.BoolFlag{
			Name:        "allowArrayToBeNull",
			Usage:       "Allow array type to be null",
			Destination: &allowArrayToBeNull,
		},
		&cli.BoolFlag{
			Name:        "allowRequiredPropertyToBeEmpty",
			Usage:       "Allow required property to be empty value, but still should be existed",
			Destination: &allowRequiredPropertyToBeEmpty,
		},
	).
	SetAction(action)

var ramlFile string
var checkRAMLVersion bool
var ignoreUnusedAnnotation bool
var ignoreUnusedTrait bool
var allowIntegerToBeNumber bool
var allowArrayToBeNull bool
var allowRequiredPropertyToBeEmpty bool

func action(c *cli.Context) (err error) {
	ramlParser := parser.NewParser()

	if err = ramlParser.Config(parserConfig.CheckRAMLVersion, checkRAMLVersion); err != nil {
		return
	}
	if err = ramlParser.Config(parserConfig.IgnoreUnusedAnnotation, ignoreUnusedAnnotation); err != nil {
		return
	}
	if err = ramlParser.Config(parserConfig.IgnoreUnusedTrait, ignoreUnusedTrait); err != nil {
		return
	}

	checkOptions := []parser.CheckValueOption{
		parser.CheckValueOptionAllowIntegerToBeNumber(allowIntegerToBeNumber),
		parser.CheckValueOptionAllowArrayToBeNull(allowArrayToBeNull),
		parser.CheckValueOptionAllowRequiredPropertyToBeEmpty(allowRequiredPropertyToBeEmpty),
	}
	if err = ramlParser.Config(parserConfig.CheckValueOptions, checkOptions); err != nil {
		return
	}

	rootdoc, err := ramlParser.ParseFile(ramlFile)
	if err != nil {
		return
	}

	jsondata, err := jsonex.MarshalIndent(rootdoc, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(jsondata))

	return
}
