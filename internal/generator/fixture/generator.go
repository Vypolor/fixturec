package fixture

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Vypolor/fixturec/internal/generator/dto"
	"github.com/Vypolor/fixturec/internal/utils"
)

const outFileName = "fixture_test.go"

type Data struct {
	ImplAlias   string
	PackageName string
	StructName  string
	Fields      []dto.FieldInfo
}

func GenerateFixtureFile(pkgDir, pkgName, structName string, fields []dto.FieldInfo) error {
	implAlias := getStructNameLower(structName)
	data := Data{PackageName: pkgName, StructName: structName, Fields: fields, ImplAlias: implAlias}

	tmpl := template.Must(template.New(templateName).Funcs(template.FuncMap{
		mockAlias: func(pkgPath string) string {
			if pkgPath == "" {
				return "local_mock"
			}

			return fmt.Sprintf("%s_%s", filepath.Base(pkgPath), "mock")
		},
		typeShort: func(full string) string {
			parts := strings.Split(full, ".")

			return fmt.Sprintf("Mock%s", parts[len(parts)-1])
		},
	}).Parse(fixtureTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("tmpl.Execute: %w", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format.Source: %w", err)
	}

	out := filepath.Join(pkgDir, outFileName)
	if err = os.WriteFile(out, src, utils.FilePerm0644); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	return nil
}

func getStructNameLower(structName string) string {
	s := structName
	if s == "" {
		return s
	}
	if s[0] >= 'A' && s[0] <= 'Z' {
		return string(s[0]+('a'-'A')) + s[1:]
	}
	return s
}
