package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/NHGRI/clu/pkg/global"
	"github.com/NHGRI/clu/pkg/input"
)

func TestRun_Hosts(t *testing.T) {
	// Save global state to restore after test
	origMockDir := global.Config.MockDir
	origRunner := input.CommandRunner
	origReader := input.FileReader
	origDebug := global.Config.Debug
	origNet := global.Config.NetEnabled

	defer func() {
		global.Config.MockDir = origMockDir
		input.CommandRunner = origRunner
		input.FileReader = origReader
		global.Config.Debug = origDebug
		global.Config.NetEnabled = origNet
	}()

	hosts := []string{"host1", "host2", "host3"}
	tiers := []int{1, 2, 3}

	for _, host := range hosts {
		for _, tier := range tiers {
			t.Run(fmt.Sprintf("%s/tier%d", host, tier), func(t *testing.T) {
				// Locate the expected output file
				// We are in cmd/, testdata is in ../testdata
				expectedFile := filepath.Join("..", "testdata", host, fmt.Sprintf("tier%d.output", tier))
				expectedBytes, err := os.ReadFile(expectedFile)
				if err != nil {
					t.Fatalf("failed to read expected output file %s: %v", expectedFile, err)
				}
				expected := string(expectedBytes)

				// EnableMockMode now automatically resolves testdata path
				args := []string{"facts", "--mock-dir", host, "--tier", fmt.Sprintf("%d", tier)}
				stdout := &bytes.Buffer{}
				stderr := &bytes.Buffer{}

				// Execute run
				exitCode := run(args, stdout, stderr)

				if exitCode != 0 {
					t.Errorf("run() failed with exit code %d. Stderr: %s", exitCode, stderr.String())
				}

				output := stdout.String()

				// Exclude clu.* lines from comparison as they vary by environment.
				re := regexp.MustCompile(`(?m)^clu\..*(\r?\n|$)`)
				output = re.ReplaceAllString(output, "")
				expected = re.ReplaceAllString(expected, "")

				// Trim trailing whitespace from lines to avoid mismatches.
				reTrim := regexp.MustCompile(`(?m)[ \t]+$`)
				output = reTrim.ReplaceAllString(output, "")
				expected = reTrim.ReplaceAllString(expected, "")

				if strings.TrimSpace(output) != strings.TrimSpace(expected) {
					showDiff(t, host, tier, output, expected)
				}
			})
		}
	}
}

func TestRun_Requires(t *testing.T) {
	// Save global state to restore after test
	origMockDir := global.Config.MockDir
	origRunner := input.CommandRunner
	origReader := input.FileReader
	origDebug := global.Config.Debug
	origNet := global.Config.NetEnabled

	defer func() {
		global.Config.MockDir = origMockDir
		input.CommandRunner = origRunner
		input.FileReader = origReader
		global.Config.Debug = origDebug
		global.Config.NetEnabled = origNet
	}()

	// Test list mode
	args := []string{"requires", "list", "--mock-dir", "host1"}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run(args, stdout, stderr)
	if exitCode != 0 {
		t.Errorf("run() failed with exit code %d. Stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Listing Requirements:") {
		t.Errorf("output missing header. Got:\n%s", output)
	}
}

func showDiff(t *testing.T, host string, tier int, got, want string) {
	t.Helper()
	gotLines := strings.Split(strings.TrimSpace(got), "\n")
	wantLines := strings.Split(strings.TrimSpace(want), "\n")

	for i := 0; i < len(gotLines) || i < len(wantLines); i++ {
		lineG := ""
		lineW := ""
		if i < len(gotLines) {
			lineG = gotLines[i]
		}
		if i < len(wantLines) {
			lineW = wantLines[i]
		}
		if lineG != lineW {
			t.Errorf("mismatch %s tier %d line %d:\n  got : %q\n  want: %q", host, tier, i+1, lineG, lineW)
			return
		}
	}
}

func TestRun_Collector(t *testing.T) {
	// Save global state
	origMockDir := global.Config.MockDir
	origRunner := input.CommandRunner
	origReader := input.FileReader
	origDebug := global.Config.Debug
	origNet := global.Config.NetEnabled

	defer func() {
		global.Config.MockDir = origMockDir
		input.CommandRunner = origRunner
		input.FileReader = origReader
		global.Config.Debug = origDebug
		global.Config.NetEnabled = origNet
	}()

	args := []string{"collector", "--mock-dir", "host1"}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	exitCode := run(args, stdout, stderr)
	if exitCode != 0 {
		t.Errorf("run() failed with exit code %d. Stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	// Expected output format: "Collection created at /tmp/clu_HOSTNAME.tgz"
	if !strings.Contains(output, "Collection created at") {
		t.Fatalf("output missing success message. Got:\n%s", output)
	}

	// Extract filename to verify existence and cleanup
	parts := strings.Fields(output)
	tarballPath := parts[len(parts)-1]

	if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
		t.Errorf("expected tarball not found at %s", tarballPath)
	} else {
		defer os.Remove(tarballPath)
	}
}
