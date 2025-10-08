package fixture

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Vypolor/fixturec/internal/generator/dto"
	"github.com/stretchr/testify/require"
)

func TestGenerateFixtureFile(t *testing.T) {
	t.Parallel()

	pkgName := "gen"
	structName := "ToGen"

	dir, err := os.MkdirTemp("", "fixturec")
	if err != nil {
		t.Fatal(err)
	}

	fields := []dto.FieldInfo{
		{
			FieldName: "myType1",
			TypeName:  "github.com/Vypolor/without_external/mypackage1.MyType1",
			PkgPath:   "github.com/Vypolor/without_external/mypackage1",
		},
		{
			FieldName: "myType2",
			TypeName:  "github.com/Vypolor/without_external/mypackage2.MyType2",
			PkgPath:   "github.com/Vypolor/without_external/mypackage2",
		},
	}

	require.NoError(t, GenerateFixtureFile(dir, pkgName, structName, fields))

	// read golden file
	golden, err := os.ReadFile(filepath.Join("testdata", "gen.golden"))
	require.NoError(t, err)

	// read generated file
	genFilePath := filepath.Join(dir, "fixture_test.go")
	generated, err := os.ReadFile(genFilePath)
	require.NoError(t, err)

	require.Equal(t, golden, generated)

	// delete temp dir
	require.NoError(t, os.RemoveAll(dir))
}
