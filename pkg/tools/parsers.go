package tools

import (
	"os/exec"
	"regexp"
	"strings"
)

// parseVersionFromOutput extracts version number from command output
func parseVersionFromOutput(output string) string {
	// Get first line
	firstLine := strings.Split(output, "\n")[0]

	// Common patterns: "name version X.Y.Z", "name X.Y.Z", "name-X.Y.Z", "name, version X.Y.Z"
	// Match version numbers like X.Y.Z, X.Y, or just X
	re := regexp.MustCompile(`\d+\.\d+(\.\d+)?`)
	matches := re.FindString(firstLine)

	if matches != "" {
		return matches
	}

	return "unknown"
}

func VersionSimpleParser(path string) string {
	// Run program --version and parse first line
	cmd := exec.Command(path, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return parseVersionFromOutput(string(output))
}

func VersionSubcmdParser(path string) string {
	// Run program version (subcommand) and parse first line
	cmd := exec.Command(path, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	return parseVersionFromOutput(string(output))
}

func VersionOpensshParser(path string) string {
	// OpenSSH uses -V (capital V, single dash) and outputs to stderr
	cmd := exec.Command(path, "-V")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// Output format: "OpenSSH_9.9p2, LibreSSL 3.3.6"
	// Extract just the version after "OpenSSH_"
	firstLine := strings.Split(string(output), "\n")[0]

	// Match OpenSSH_X.YpZ (p is literal, followed by patch number)
	re := regexp.MustCompile(`OpenSSH_(\d+\.\d+p\d+)`)
	matches := re.FindStringSubmatch(firstLine)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}

func VersionPassParser(path string) string {
	// pass has unique ASCII art output with version on line 4
	cmd := exec.Command(path, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// Output format has version on line 4 like "=                  v1.7.4                  ="
	lines := strings.Split(string(output), "\n")
	if len(lines) < 4 {
		return "unknown"
	}

	// Extract version from line with "v1.7.4" format
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(lines[3])

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}
