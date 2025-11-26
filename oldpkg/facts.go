package pkg

type tier int

const (
	Unknown tier = iota
	One
	Two
	Three
)

func (t tier) String() string {
	switch t {
	case Unknown:
		return "UNKNOWN"
	case One:
		return "TIER1"
	case Two:
		return "TIER2"
	case Three:
		return "TIER3"
	default:
		return "INVALID TIER"
	}
}

type Facts map[string]string

// var f Facts f = make(Facts)
// f := Facts{}
func NewFacts() Facts {
	return make(Facts)
}

func (f Facts) Add(name string, value string) {
	f[name] = value
}
