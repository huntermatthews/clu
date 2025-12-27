package sources

// Go port of src/clu/sources/no_salt.py
// Detects presence of a sentinel file /no_salt and records existence and reason.

import (
	"strings"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// NoSalt reports whether the /no_salt file exists and its content as a reason.
type NoSalt struct{}

// Provides registers fact keys produced by this source.
func (n *NoSalt) Provides(p types.Provides) {
	p["salt.no_salt.exists"] = n
	p["salt.no_salt.reason"] = n
}

// Requires declares file dependency.
func (n *NoSalt) Requires(r *types.Requires) {
	r.Files = append(r.Files, "/no_salt")
}

// Parse reads /no_salt optionally. If missing or empty sets exists False and returns.
// Otherwise sets exists True and reason to trimmed file content.
func (n *NoSalt) Parse(f *types.Facts) {
	// Use DI FileReader (returns content,error). Treat any error or empty
	// trimmed content as missing (optional semantics).
	data, err := pkg.FileReader("/no_salt")
	if err != nil || strings.TrimSpace(data) == "" {
		f.Set("salt.no_salt.exists", "False")
		return
	}
	f.Set("salt.no_salt.exists", "True")
	f.Set("salt.no_salt.reason", strings.TrimSpace(data))
}
