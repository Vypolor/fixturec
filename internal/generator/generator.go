package generator

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/Vypolor/fixturec/internal/generator/dto"
	"github.com/Vypolor/fixturec/internal/generator/fixture"
	"github.com/Vypolor/fixturec/internal/generator/mock"
	"github.com/Vypolor/fixturec/internal/utils"
	"golang.org/x/tools/go/packages"
)

type Generator struct {
	dir      string
	typeName string
}

func NewGenerator(dir, typeName string) *Generator {
	return &Generator{dir: dir, typeName: typeName}
}

func (g *Generator) Generate() error {
	pkg, pkgDir, err := utils.LoadPackageByPattern(g.dir)
	if err != nil {
		return fmt.Errorf("utils.LoadPackageByPattern: %w", err)
	}

	fields, err := g.getInterfacesStructFields(pkg)
	if err != nil {
		return fmt.Errorf("getInterfacesStructFields: %w", err)
	}

	if len(fields) == 0 {
		log.Printf("no interfaces fields found for struct type: %s", g.typeName)

		return nil
	}

	if err = mock.GenerateMocks(pkgDir, fields); err != nil {
		return fmt.Errorf("mock.GenerateMocks: %w", err)
	}

	// reload package
	pkg, pkgDir, err = utils.LoadPackageByPattern(g.dir)
	if err != nil {
		return fmt.Errorf("utils.LoadPackageByPattern: %w", err)
	}

	if err = fixture.GenerateFixtureFile(pkgDir, pkg.Name, g.typeName, fields); err != nil {
		return fmt.Errorf("fixture.GenerateFixtureFile: %w", err)
	}

	return nil
}

func (g *Generator) getInterfacesStructFields(pkg *packages.Package) ([]dto.FieldInfo, error) {
	obj := pkg.Types.Scope().Lookup(g.typeName)
	if obj == nil {
		return nil, fmt.Errorf("type %s not found", g.typeName)
	}
	tn, ok := obj.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("%s is not a type", g.typeName)
	}
	strct, ok := tn.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s is not a struct", g.typeName)
	}

	modulePath := ""
	if pkg.Module != nil {
		modulePath = pkg.Module.Path
	}

	res := make([]dto.FieldInfo, 0, strct.NumFields())
	for i := 0; i < strct.NumFields(); i++ {
		field := strct.Field(i)
		fieldType := field.Type()

		if _, ok := fieldType.Underlying().(*types.Interface); ok {
			pkgPath := ""
			if named, ok := fieldType.(*types.Named); ok {
				if named.Obj() != nil && named.Obj().Pkg() != nil {
					pkgPath = named.Obj().Pkg().Path()
				}
			}
			if pkgPath == "" || strings.HasPrefix(pkgPath, modulePath) {
				res = append(res, dto.FieldInfo{
					FieldName: field.Name(),
					TypeName:  types.TypeString(field.Type(), func(p *types.Package) string { return p.Path() }),
					PkgPath:   pkgPath,
				})
			}
		}
	}

	return res, nil
}
