// This program generates files used to access the UAA API. It can be invoked
// by running go generate

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	uaa "github.com/cloudfoundry-community/go-uaa"
)

var typesToProcess = []interface{}{
	uaa.Client{},
	uaa.Group{},
	uaa.User{},
	uaa.IdentityZone{},
	uaa.MFAProvider{},
}

func main() {
	for _, typ := range typesToProcess {
		rtype := reflect.TypeOf(typ)
		if rtype.Kind() != reflect.Struct {
			log.Printf("[warn]: %s is not a struct type...skipping\n", rtype.Name())
			continue
		}

		typeName := rtype.Name()
		pluralTypeName := typeName + "s"

		t := typeGenerator{
			ModelTypeName:       typeName,
			ModelPluralTypeName: pluralTypeName,
			IDFieldName:         "ID",
			SupportsAttributes:  true,
			SupportsPaging:      true,
		}
		if typeName == "Client" {
			t.SupportsAttributes = false
		}

		if typeName == "IdentityZone" || typeName == "MFAProvider" {
			t.SupportsPaging = false
		}

		for i := 0; i < rtype.NumField(); i++ {
			field := rtype.Field(i)

			fieldTag := field.Tag.Get("generator")
			if fieldTag == "id" {
				t.IDFieldName = field.Name
			}
		}

		sourceBuf := &bytes.Buffer{}
		specBuf := &bytes.Buffer{}

		generate(sourceBuf, specBuf, t)

		if sourceBuf.Len() > 0 {
			writeFile(sourceBuf.Bytes(), strings.ToLower(fmt.Sprintf("generated_%s.go", t.ModelTypeName)))
		}
		if specBuf.Len() > 0 {
			writeFile(specBuf.Bytes(), strings.ToLower(fmt.Sprintf("generated_%s_test.go", t.ModelTypeName)))
		}
	}
}

func generate(sourceBuf io.Writer, specBuf io.Writer, t typeGenerator) {
	if err := modelTmpl.Execute(sourceBuf, t); err != nil {
		log.Fatalf("generating model code: %v", err)
	}

	if err := specTmpl.Execute(specBuf, t); err != nil {
		log.Fatalf("generating test code: %v", err)
	}
}

func writeFile(buf []byte, filename string) {
	src, err := format.Source(buf)
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		src = buf
	}

	if err = ioutil.WriteFile(filename, src, 0644); err != nil {
		log.Fatalf("writing output [%s]: %v", filename, err)
	}
}

func tolower(s string) string {
	return strings.ToLower(s)
}

func join(s ...string) string {
	return strings.Join(s, ", ")
}

func joinfields(fields []structField, includeMetadata bool) string {
	names := make([]string, 0, len(fields))
	for i := range fields {
		names = append(names, tolower(fields[i].Name))
	}
	if includeMetadata {
		names = append(names, "metadata")
	}

	return join(names...)
}

func joinfieldsprefixlower(fields []structField, prefix string) string {
	return tolower(joinfieldsprefix(fields, prefix))
}

func joinfieldsprefix(fields []structField, prefix string) string {
	names := make([]string, 0, len(fields))
	for i := range fields {
		names = append(names, prefix+fields[i].Name)
	}

	return join(names...)
}

type structField struct {
	ModelTypeName       string // name of the parent model object
	ModelPluralTypeName string // plural name of the parent model object
	Name                string // name of the field
}

type typeGenerator struct {
	ModelTypeName       string        // the name of the model type
	ModelPluralTypeName string        // the plural of the model type
	IDFieldName         string        // the field name for the ID
	SupportsAttributes  bool          // attributes can be supplied when listing
	SupportsPaging      bool          // paging is supported
	Fields              []structField // fields on the struct we're generating for (converted to columns)
}

var (
	mode string
)

var modelTmpl = template.Must(template.New("modelTempl").Funcs(template.FuncMap{
	"tolower": tolower,
}).Parse(load("model.gotemplate")))

var specTmpl = template.Must(template.New("specTempl").Funcs(template.FuncMap{
	"tolower": tolower,
}).Parse(load("spec.gotemplate")))

func load(name string) string {
	_, file, _, _ := runtime.Caller(1)
	b, err := ioutil.ReadFile(filepath.Join(filepath.Dir(file), name))
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
