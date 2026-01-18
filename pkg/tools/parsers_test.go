package tools

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestVersionSimpleParser_Python3(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "python3_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "3.14.0"

	if result != expected {
		t.Errorf("Python3: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Jq(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "jq_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "1.8.1"

	if result != expected {
		t.Errorf("Jq: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Git(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "git_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "2.51.2"

	if result != expected {
		t.Errorf("Git: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Fish(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "fish_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "4.1.2"

	if result != expected {
		t.Errorf("Fish: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Just(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "just_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "1.43.0"

	if result != expected {
		t.Errorf("Just: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Bash(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "bash_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "5.3.3"

	if result != expected {
		t.Errorf("Bash: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Curl(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "curl_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "8.7.1"

	if result != expected {
		t.Errorf("Curl: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Gpg(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "gpg_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "2.4.8"

	if result != expected {
		t.Errorf("Gpg: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Make(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "make_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "4.4.1"

	if result != expected {
		t.Errorf("Make: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Pip3(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "pip3_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "25.3"

	if result != expected {
		t.Errorf("Pip3: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Pipx(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "pipx_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "1.8.0"

	if result != expected {
		t.Errorf("Pipx: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Uv(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "uv_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "0.9.7"

	if result != expected {
		t.Errorf("Uv: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Zsh(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "zsh_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "5.9"

	if result != expected {
		t.Errorf("Zsh: expected %q, got %q", expected, result)
	}
}

func TestVersionSimpleParser_Awk(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "awk_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	// awk version 20200816 - this won't match the X.Y pattern, so expect "unknown"
	expected := "unknown"

	if result != expected {
		t.Errorf("Awk: expected %q, got %q", expected, result)
	}
}

func TestVersionOpensshParser(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "ssh_-V")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Parse using OpenSSH-specific logic
	firstLine := strings.Split(string(content), "\n")[0]
	re := regexp.MustCompile(`OpenSSH_(\d+\.\d+p\d+)`)
	matches := re.FindStringSubmatch(firstLine)

	var result string
	if len(matches) > 1 {
		result = matches[1]
	} else {
		result = "unknown"
	}

	expected := "9.9p2"

	if result != expected {
		t.Errorf("OpenSSH: expected %q, got %q", expected, result)
	}
}

func TestVersionPassParser(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "pass_--version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Parse using pass-specific logic
	lines := strings.Split(string(content), "\n")
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)
	var result string
	if len(lines) >= 4 {
		matches := re.FindStringSubmatch(lines[3])
		if len(matches) > 1 {
			result = matches[1]
		} else {
			result = "unknown"
		}
	} else {
		result = "unknown"
	}

	expected := "1.7.4"

	if result != expected {
		t.Errorf("Pass: expected %q, got %q", expected, result)
	}
}

func TestVersionSubcmdParser_Go(t *testing.T) {
	testFile := filepath.Join("..", "..", "testdata", "host_tools", "_programs", "go_version")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	result := parseVersionFromOutput(string(content))
	expected := "1.25.3"

	if result != expected {
		t.Errorf("Go: expected %q, got %q", expected, result)
	}
}
