package pkg

// Requires mirrors the Python dataclass in requires.py, holding lists of
// files, programs, APIs, and facts a source depends on.
// In Python: dataclass with update() that extends each list.
// In Go we implement Update with slice appends.

type Requires struct {
	Files    []string
	Programs []string
	APIs     []string
	Facts    []string
}

// NewRequires constructs an empty Requires.
func NewRequires() *Requires {
	return &Requires{}
}

// Update extends the receiver's slices with those from other, returning the receiver.
func (r *Requires) Update(other *Requires) *Requires {
	if other == nil {
		return r
	}
	r.Files = append(r.Files, other.Files...)
	r.Programs = append(r.Programs, other.Programs...)
	r.APIs = append(r.APIs, other.APIs...)
	r.Facts = append(r.Facts, other.Facts...)
	return r
}
