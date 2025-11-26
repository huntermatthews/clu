package pkg

// Port of src/clu/__about__.py providing package metadata.
// In Python these values are derived dynamically from distribution metadata.
// In Go we expose a set of constants/variables that can be overridden at build time
// via -ldflags (e.g., -ldflags "-X github.com/huntermatthews/clu/pkg.Version=1.2.3").
// Adjust as needed if integrating with module-aware versioning.

var (
	Title       = "clu"
	Version     = "0.dev0+unknown" // override at build time
	Summary     = "clu utility"    // manual copy; Python pulled from metadata
	License     = "UNKNOWN"        // set appropriately
	AuthorEmail = "UNKNOWN"        // set appropriately

	// MinimumPython mirrored from Python metadata; here just informational.
	MinimumPython      = "3.11"
	MinimumPythonMajor = 3
	MinimumPythonMinor = 11
)

// SetVersion allows runtime override if desired (e.g., from CLI flag or env).
func SetVersion(v string) { Version = v }

// Metadata returns a snapshot map of the metadata (convenience helper).
func Metadata() map[string]string {
	return map[string]string{
		"title":        Title,
		"version":      Version,
		"summary":      Summary,
		"license":      License,
		"author_email": AuthorEmail,
		"min_python":   MinimumPython,
	}
}
