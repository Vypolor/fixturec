package mock

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Vypolor/fixturec/internal/generator/dto"
	"github.com/Vypolor/fixturec/internal/utils"
	"golang.org/x/tools/go/packages"
)

func GenerateMocks(pkgDir string, fields []dto.FieldInfo) error {
	pkgMap := map[string][]dto.FieldInfo{}
	for _, field := range fields {
		pkgMap[field.PkgPath] = append(pkgMap[field.PkgPath], field)
	}

	for pkgPath, fis := range pkgMap {
		pDir, err := getPackageDir(pkgDir, pkgPath)
		if err != nil {
			return fmt.Errorf("getPackageDir: %w", err)
		}

		cfg := &packages.Config{
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
				packages.NeedTypes | packages.NeedTypesInfo,
			Dir: pDir,
			// Env: append(os.Environ(), utils.EnvGoSumDbDefaultValue, utils.EnvGoSumDbDefaultValue),
		}
		pkgs, err := packages.Load(cfg, "./")
		if err != nil || len(pkgs) == 0 {
			return fmt.Errorf("failed to load pkg %s", pDir)
		}
		targetPkg := pkgs[0]

		filesToModify := map[string]bool{}
		for _, fi := range fis {
			parts := strings.Split(fi.TypeName, ".")
			typeName := parts[len(parts)-1]

			obj := targetPkg.Types.Scope().Lookup(typeName)
			if obj == nil {
				continue
			}
			posn := targetPkg.Fset.Position(obj.Pos())
			if posn.Filename != "" {
				filesToModify[posn.Filename] = true
			}
		}

		for filename := range filesToModify {
			b, _ := os.ReadFile(filename)
			if strings.Contains(string(b), "//go:generate mockgen") {
				continue
			}
			dirLine := fmt.Sprintf("//go:generate mockgen -destination=mock/mock_gen.go -package=mock -source=./%s\n",
				filepath.Base(filename))

			if err = os.WriteFile(filename, append([]byte(dirLine), b...), utils.FilePerm0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", filename, err)
			}
			fmt.Printf("Inserted mockgen directive into %s\n", filename)
		}

		cmd := exec.Command("go", "generate", "./...")
		cmd.Dir = pDir
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("cmd.Run: %w", err)
		}
	}

	return nil
}

func getPackageDir(startDir, pkgPath string) (string, error) {
	if pkgPath == "" {
		return startDir, nil
	}
	if fi, err := os.Stat(pkgPath); err == nil && fi.IsDir() {
		return pkgPath, nil
	}
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedModule,
		Dir:  startDir,
		// Env:  append(os.Environ(), utils.EnvGoSumDbDefaultValue, utils.EnvGoProxyDirect),
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return "", err
	}
	if packages.PrintErrors(pkgs) > 0 || len(pkgs) == 0 {
		return "", fmt.Errorf("package %s not found", pkgPath)
	}
	first := pkgs[0]
	if len(first.GoFiles) == 0 {
		return "", fmt.Errorf("could not determine directory for package %s", pkgPath)

	}

	return filepath.Dir(first.GoFiles[0]), nil
}
