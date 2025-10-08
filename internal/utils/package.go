package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName |
	packages.NeedFiles | packages.NeedSyntax |
	packages.NeedTypes | packages.NeedTypesInfo |
	packages.NeedModule

const (
	EnvGoSumDbDefaultValue = "GOSUMDB='sum.golang.org'"
)

func LoadPackageByPattern(pattern string) (*packages.Package, string, error) {
	if fi, err := os.Stat(pattern); err == nil && fi.IsDir() {
		pkgs, err := loadPackage(pattern)
		if err != nil {
			return nil, "", fmt.Errorf("loadPackage: %w", err)
		}

		return pkgs[0], pattern, nil
	}

	pkgs, err := loadPackage(pattern)
	if err != nil {
		return nil, "", fmt.Errorf("loadPackage: %w", err)
	}

	var (
		pkg    = pkgs[0]
		pkgDir string
	)
	if len(pkg.GoFiles) > 0 {
		pkgDir = filepath.Dir(pkg.GoFiles[0])
	}

	return pkg, pkgDir, nil
}

func loadPackage(pattern string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: loadMode,
		Dir:  pattern,
		// Env:  append(os.Environ(), EnvGoSumDbDefaultValue),
	}
	pkgs, err := packages.Load(cfg, "./")
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 || len(pkgs) == 0 {
		return nil, fmt.Errorf("failed to load package from %s", pattern)
	}

	return pkgs, nil
}
