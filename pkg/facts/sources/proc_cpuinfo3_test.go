package sources

import (
	"errors"
	"testing"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

func TestProcCpuinfo3Provides(t *testing.T) {
	src := &ProcCpuinfo3{}
	p := types.Provides{}
	src.Provides(p)
	for k := range procCpuinfo3Facts {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

// Build a cpuinfo block helper
func cpuBlock(entries map[string]string) string {
	var sb string
	for k, v := range entries {
		sb += k + "\t: " + v + "\n"
	}
	return sb
}

func TestProcCpuinfo3ParseCases(t *testing.T) {
	orig := input.FileReader
	defer func() { input.FileReader = orig }()

	// Case: single logical CPU
	single := cpuBlock(map[string]string{
		"processor":   "0",
		"vendor_id":   "GenuineIntel",
		"model name":  "Intel(R) Test CPU",
		"physical id": "0",
		"core id":     "0",
		"cpu cores":   "1",
	})

	// Case: two sockets, two cores each (4 unique core pairs), 4 logical
	multi := cpuBlock(map[string]string{"processor": "0", "vendor_id": "ACME", "model name": "ACME CPU", "physical id": "0", "core id": "0"}) + "\n" +
		cpuBlock(map[string]string{"processor": "1", "physical id": "0", "core id": "1"}) + "\n" +
		cpuBlock(map[string]string{"processor": "2", "physical id": "1", "core id": "0"}) + "\n" +
		cpuBlock(map[string]string{"processor": "3", "physical id": "1", "core id": "1"})

	// Case: physical id present but no core id; cpu cores present to allow fallback
	fallback := cpuBlock(map[string]string{"processor": "0", "vendor_id": "V", "model name": "M", "physical id": "0", "cpu cores": "2"}) + "\n" +
		cpuBlock(map[string]string{"processor": "1", "physical id": "0"}) + "\n" +
		cpuBlock(map[string]string{"processor": "2", "physical id": "1"}) + "\n" +
		cpuBlock(map[string]string{"processor": "3", "physical id": "1"})

	cases := []struct {
		name string
		data string
		want map[string]string
	}{
		{"single", single, map[string]string{"phy.cpu.vendor": "GenuineIntel", "phy.cpu.model": "Intel(R) Test CPU", "phy.cpu.threads": "1", "phy.cpu.sockets": "1", "phy.cpu.cores": "1"}},
		{"multi", multi, map[string]string{"phy.cpu.vendor": "ACME", "phy.cpu.model": "ACME CPU", "phy.cpu.threads": "4", "phy.cpu.sockets": "2", "phy.cpu.cores": "4"}},
		{"fallback", fallback, map[string]string{"phy.cpu.vendor": "V", "phy.cpu.model": "M", "phy.cpu.threads": "4", "phy.cpu.sockets": "2", "phy.cpu.cores": "4"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input.FileReader = func(path string) (string, error) { return c.data, nil }
			f := types.NewFactDB()
			src := &ProcCpuinfo3{}
			src.Parse(f)
			for k, want := range c.want {
				got, _ := f.Get(k)
				if got != want {
					t.Fatalf("%s: want %s=%s got %s", c.name, k, want, got)
				}
			}
		})
	}
}

func TestProcCpuinfo3MissingFile(t *testing.T) {
	orig := input.FileReader
	input.FileReader = func(path string) (string, error) { return "", errors.New("missing") }
	defer func() { input.FileReader = orig }()

	f := types.NewFactDB()
	src := &ProcCpuinfo3{}
	src.Parse(f)

	// cores and threads should be ParseFailMsg; vendor/model empty; sockets empty
	if got, _ := f.Get("phy.cpu.cores"); got != types.ParseFailMsg {
		t.Fatalf("expected cores ParseFailMsg got %s", got)
	}
	if got, _ := f.Get("phy.cpu.threads"); got != types.ParseFailMsg {
		t.Fatalf("expected threads ParseFailMsg got %s", got)
	}
	if got, _ := f.Get("phy.cpu.vendor"); got != "" {
		t.Fatalf("expected empty vendor got %s", got)
	}
	if got, _ := f.Get("phy.cpu.model"); got != "" {
		t.Fatalf("expected empty model got %s", got)
	}
	if got, _ := f.Get("phy.cpu.sockets"); got != "" {
		t.Fatalf("expected empty sockets got %s", got)
	}
}
