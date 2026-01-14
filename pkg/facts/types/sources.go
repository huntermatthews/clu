package types

const (
	ParseFailMsg   = "Unknown/Error"
	NetDisabledMsg = "Unknown - Network Queries Disabled"
)

// Sources is the common interface for fact sources.
// Methods mirror Python abstract methods: provides, requires, parse.
type Sources interface {
	Provides(p Provides)
	Requires(r *Requires)
	Parse(f *FactDB)
}
