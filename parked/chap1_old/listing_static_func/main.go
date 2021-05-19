package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"text/tabwriter"
)

type Output struct {
	File      string `json:"file"`
	Name      string `json:"name"`
	NumInput  int    `json:"num_input"`
	NumReturn int    `json:"num_return"`
}

func main() {

	var files []*ast.File

	var formatString, outputFormat string
	flag.StringVar(&formatString, "f", "", "format string")
	flag.StringVar(&outputFormat, "o", "table", "output format")
	flag.Parse()

	if outputFormat != "table" && len(formatString) != 0 {
		log.Fatal("Format string only valid with table output format")
	}

	fset := token.NewFileSet()
	for _, goFile := range flag.Args() {
		f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, f)
	}
	w := os.Stdout
	var functions []Output

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			var fName *ast.Ident
			var fType *ast.FuncType
			var fPos token.Pos
			switch x := n.(type) {
			case *ast.FuncDecl:
				fName = x.Name
				fType = x.Type
				fPos = fType.Pos()
			}
			if fName != nil {
				var nParams, nResults int
				if fType.Params != nil {
					nParams = len(fType.Params.List)
				}
				if fType.Results != nil {
					nResults = len(fType.Results.List)
				}
				p := fset.Position(fPos)
				o := Output{
					File:      p.Filename,
					Name:      fmt.Sprintf("%v", fName),
					NumInput:  nParams,
					NumReturn: nResults,
				}
				functions = append(functions, o)
			}
			return true
		})
	}
	if outputFormat == "json" {
		json, err := json.MarshalIndent(functions, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(json))
	} else {
		const padding = 5
		w := tabwriter.NewWriter(os.Stdout, 8, 16, padding, ' ', 0)
		if len(formatString) == 0 {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "File", "Name", "NumInput", "NumReturn")
			formatString = "{{.File}}\t{{.Name}}\t{{.NumInput}}\t{{.NumReturn}}"
		}
		tmpl := template.New("test")
		tmpl, err := tmpl.Parse(formatString)
		if err != nil {
			log.Fatal("Error Parsing template: ", err)
			return
		}
		for _, f := range functions {
			err1 := tmpl.Execute(w, f)
			if err1 != nil {
				log.Fatal("Error executing template: ", err1)
			}
			fmt.Fprintln(w)
		}
		w.Flush()
	}
}
