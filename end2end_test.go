package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	testWithoutExternal = "without_external"

	testDataDir = "testdata"
)

var GOEXE = ""

func init() {
	if runtime.GOOS == "windows" {
		GOEXE = ".exe"
	}
}

func TestEndToEnd(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "fixturec")
	if err != nil {
		t.Fatalf("os.MkdirTemp: %v", err)
	}

	// Create fixturec in temporary directory
	fixturec := filepath.Join(dir, fmt.Sprintf("fixturec%s", GOEXE))
	if err = execCommand("go", "build", "-o", fixturec); err != nil {
		t.Fatalf("execCommand: %v", err)
	}

	// Read testdata directory
	fd, err := os.Open(testDataDir)
	if err != nil {
		t.Fatalf("os.Open: %v", err)
	}
	defer fd.Close()

	testNames, err := fd.Readdirnames(-1)
	if err != nil {
		t.Fatalf("Readdirnames: %v", err)
	}

	// Generate and run tests
	for _, testName := range testNames {
		var (
			genPackage    string
			genStructName string
		)
		switch testName {
		case testWithoutExternal:
			genPackage = "gen"
			genStructName = "ToGen"
		default:
			t.Fatalf("unrecognized test name %q", testName)
		}

		e2eTest(t, fixturec, testName, genPackage, genStructName)
	}
}

func e2eTest(t *testing.T, fixturec, testName, genPackage, genStructName string) {
	t.Logf("running test %s", testName)

	testRoot := filepath.Join(testDataDir, testName)
	genPackagePath := filepath.Join(testRoot, genPackage)

	callFixturec(t, fixturec, genPackagePath, genStructName)

	runTests(t, testRoot)
}

func callFixturec(t *testing.T, fixturec, genPackage, genStructName string) {
	args := []string{"-t", genStructName}

	if err := execCommandInDir(genPackage, fixturec, args...); err != nil {
		t.Fatalf("execCommandInDir: %v", err)
	}
}

func runTests(t *testing.T, dir string) {
	if err := execCommandInDir(dir, "go", "test", "./...", "-v"); err != nil {
		t.Fatalf("execCommandInDir: %v", err)
	}
}

func execCommand(name string, args ...string) error {
	return execCommandInDir(".", name, args...)
}

func execCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
