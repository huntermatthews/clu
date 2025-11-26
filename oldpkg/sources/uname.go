package sources

import (
	"log"
	"os/exec"
	"strings"
)

// PARSE_FAIL_MSG is used when parsing the uname output fails.
const PARSE_FAIL_MSG = "PARSE_FAIL"

// Minimal types to match the Python source expectations in this package.
// Other code in the repo can use these or replace them with richer
// implementations later.
type Source interface{}

type Provides map[string]Source

type Requires struct {
	Programs []string
}

// Uname collects simple uname-derived facts by running `uname -snrm`.
type Uname struct{}

var unameKeys = []string{
	"os.kernel.name",
	"os.hostname",
	"os.kernel.version",
	"phy.arch",
}

// Provides registers which keys this source provides.
func (u *Uname) Provides(p Provides) {
	for _, k := range unameKeys {
		p[k] = u
	}
}

// Requires declares external programs this source depends on.
func (u *Uname) Requires(r *Requires) {
	r.Programs = append(r.Programs, "uname -snrm")
}

// textProgram runs a small text program (command with args) and returns
// its stdout and the exit code. On non-zero exit we'll return any captured
// stdout and the exit code. If the command couldn't be started we return
// empty string and -1.
func textProgram(cmdline string) (string, int) {
	parts := strings.Fields(cmdline)
	if len(parts) == 0 {
		return "", -1
	}
	name := parts[0]
	args := parts[1:]
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			// Return combined output (stdout+stderr) and exit code
			return string(out), ee.ExitCode()
		}
		log.Printf("failed to run %q: %v", cmdline, err)
		return "", -1
	}
	return string(out), 0
}

// Parse populates the provided Facts map with values from `uname -snrm`.
// If the primary key already exists this function is a no-op.
func (u *Uname) Parse(f Facts) {
	if _, ok := f["os.kernel.name"]; ok {
		return
	}

	data, rc := textProgram("uname -snrm")
	if data == "" || rc != 0 {
		for _, k := range unameKeys {
			f[k] = PARSE_FAIL_MSG
		}
		return
	}

	fields := strings.Fields(strings.TrimSpace(data))
	for i, k := range unameKeys {
		if i < len(fields) {
			f[k] = fields[i]
		} else {
			f[k] = PARSE_FAIL_MSG
		}
	}
}
