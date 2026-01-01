package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/huntermatthews/clu/pkg/global"
	"github.com/huntermatthews/clu/pkg/input"
)

func setupTest(t *testing.T) (string, func()) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	testDataDir := filepath.Join(wd, "testdata")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("failed to create testdata dir: %v", err)
	}

	origMockDir := global.Config.MockDir
	origRunner := input.CommandRunner
	origReader := input.FileReader

	return testDataDir, func() {
		global.Config.MockDir = origMockDir
		input.CommandRunner = origRunner
		input.FileReader = origReader
	}
}

func TestEnableMockMode(t *testing.T) {
	testDataDir, cleanup := setupTest(t)
	defer cleanup()

	// 1. Valid directory
	validDirName := "valid_mock"
	validDirPath := filepath.Join(testDataDir, validDirName)
	if err := os.MkdirAll(validDirPath, 0755); err != nil {
		t.Fatalf("failed to create valid mock dir: %v", err)
	}
	defer os.RemoveAll(validDirPath)

	// 2. File pretending to be directory
	fileName := "not_a_dir"
	filePath := filepath.Join(testDataDir, fileName)
	if err := os.WriteFile(filePath, []byte("fake"), 0644); err != nil {
		t.Fatalf("failed to create fake file: %v", err)
	}
	defer os.Remove(filePath)

	tests := []struct {
		name        string
		dirArg      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid Directory",
			dirArg:  validDirName,
			wantErr: false,
		},
		{
			name:        "Non-existent Directory",
			dirArg:      "does_not_exist",
			wantErr:     true,
			errContains: "mock directory not found",
		},
		{
			name:        "File instead of Directory",
			dirArg:      fileName,
			wantErr:     true,
			errContains: "mock path is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnableMockMode(tt.dirArg)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnableMockMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("EnableMockMode() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				// Verify success state
				expected := filepath.Join(testDataDir, tt.dirArg)
				if global.Config.MockDir != expected {
					t.Errorf("global.Config.MockDir = %q, want %q", global.Config.MockDir, expected)
				}

				// Verify function swaps
				if reflect.ValueOf(input.CommandRunner).Pointer() != reflect.ValueOf(input.MockTextProgram).Pointer() {
					t.Error("input.CommandRunner was not swapped to input.MockTextProgram")
				}
				if reflect.ValueOf(input.FileReader).Pointer() != reflect.ValueOf(input.MockTextFile).Pointer() {
					t.Error("input.FileReader was not swapped to input.MockTextFile")
				}
			}
		})
	}
}

func TestMockFileReading(t *testing.T) {
	testDataDir, cleanup := setupTest(t)
	defer cleanup()

	mockDirName := "read_test_mock"
	mockDirPath := filepath.Join(testDataDir, mockDirName)

	// Create mock directory structure: testdata/read_test_mock/etc
	targetDir := filepath.Join(mockDirPath, "etc")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("failed to create mock dirs: %v", err)
	}
	defer os.RemoveAll(mockDirPath)

	// Create a dummy file: testdata/read_test_mock/etc/os-release
	expectedContent := "NAME=MockOS"
	filePath := filepath.Join(targetDir, "os-release")
	if err := os.WriteFile(filePath, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("failed to write mock file: %v", err)
	}

	// Enable Mock Mode
	if err := EnableMockMode(mockDirName); err != nil {
		t.Fatalf("EnableMockMode failed: %v", err)
	}

	// Test reading the file via the swapped function
	content, err := input.FileReader("/etc/os-release")
	if err != nil {
		t.Fatalf("input.FileReader failed: %v", err)
	}

	if content != expectedContent {
		t.Errorf("got content %q, want %q", content, expectedContent)
	}
}

func TestMockCommandExecution(t *testing.T) {
	testDataDir, cleanup := setupTest(t)
	defer cleanup()

	mockDirName := "cmd_test_mock"
	mockDirPath := filepath.Join(testDataDir, mockDirName)

	// Create mock directory structure: testdata/cmd_test_mock/_programs
	programsDir := filepath.Join(mockDirPath, "_programs")
	if err := os.MkdirAll(programsDir, 0755); err != nil {
		t.Fatalf("failed to create mock dirs: %v", err)
	}
	defer os.RemoveAll(mockDirPath)

	// Define command and expected output
	cmd := "uname -r"
	expectedOutput := "5.10.0-fake\n"

	// Calculate filename used by mock system
	cmdFilename, _ := input.TransformCmdlineToFilename(cmd)
	outPath := filepath.Join(programsDir, cmdFilename)

	if err := os.WriteFile(outPath, []byte(expectedOutput), 0644); err != nil {
		t.Fatalf("failed to write mock cmd file: %v", err)
	}

	// Enable Mock Mode
	if err := EnableMockMode(mockDirName); err != nil {
		t.Fatalf("EnableMockMode failed: %v", err)
	}

	// Test execution
	out, rc, err := input.CommandRunner(cmd)
	if err != nil {
		t.Fatalf("input.CommandRunner failed: %v", err)
	}

	if rc != 0 {
		t.Errorf("got rc %d, want 0", rc)
	}
	if out != expectedOutput {
		t.Errorf("got output %q, want %q", out, expectedOutput)
	}
}
