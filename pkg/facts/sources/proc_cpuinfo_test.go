package sources

import (
	"errors"
	"testing"

	"github.com/NHGRI/clu/pkg/facts/types"
	"github.com/NHGRI/clu/pkg/input"
)

// Minimal synthetic /proc/cpuinfo fragments for testing flag progression.
var cpuinfoV1 = "processor\t: 0\nflags\t: lm cmov cx8 fpu fxsr mmx syscall sse2"
var cpuinfoV2 = "processor\t: 0\nflags\t: lm cmov cx8 fpu fxsr mmx syscall sse2 cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3"
var cpuinfoV3 = "processor\t: 0\nflags\t: lm cmov cx8 fpu fxsr mmx syscall sse2 cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3 avx avx2 bmi1 bmi2 f16c fma abm movbe xsave"
var cpuinfoV4 = "processor\t: 0\nflags\t: lm cmov cx8 fpu fxsr mmx syscall sse2 cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3 avx avx2 bmi1 bmi2 f16c fma abm movbe xsave avx512f avx512bw avx512cd avx512dq avx512vl"

func TestProcCpuinfoProvides(t *testing.T) {
	src := &ProcCpuinfo{}
	p := types.Provides{}
	src.Provides(p)
	if _, ok := p["phy.cpu.arch_version"]; !ok {
		t.Fatalf("missing provides key phy.cpu.arch_version")
	}
}

func TestProcCpuinfoVersions(t *testing.T) {
	cases := []struct {
		name string
		data string
		want string
	}{
		{"v0", "processor: 0\nflags:\t ", "x86_64_v0"},
		{"v1", cpuinfoV1, "x86_64_v1"},
		{"v2", cpuinfoV2, "x86_64_v2"},
		{"v3", cpuinfoV3, "x86_64_v3"},
		{"v4", cpuinfoV4, "x86_64_v4"},
	}
	orig := input.FileReader
	defer func() { input.FileReader = orig }()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input.FileReader = func(path string) (string, error) { return c.data, nil }
			f := types.NewFactDB()
			f.AddFact(types.Fact{Name: "phy.arch", Value: "x86_64"})
			src := &ProcCpuinfo{}
			src.Parse(f)
			got, _ := f.Get("phy.cpu.arch_version")
			if got != c.want {
				t.Fatalf("want %s got %s", c.want, got)
			}
		})
	}
}

func TestProcCpuinfoNonX86Skip(t *testing.T) {
	orig := input.FileReader
	input.FileReader = func(path string) (string, error) { return cpuinfoV4, nil }
	defer func() { input.FileReader = orig }()
	f := types.NewFactDB()
	f.AddFact(types.Fact{Name: "phy.arch", Value: "arm64"})
	src := &ProcCpuinfo{}
	src.Parse(f)
	if _, ok := f.Get("phy.cpu.arch_version"); ok {
		t.Fatalf("expected no arch_version fact on non-x86 arch")
	}
}

func TestProcCpuinfoMissingFile(t *testing.T) {
	orig := input.FileReader
	input.FileReader = func(path string) (string, error) { return "", errors.New("missing") }
	defer func() { input.FileReader = orig }()
	f := types.NewFactDB()
	f.AddFact(types.Fact{Name: "phy.arch", Value: "x86_64"})
	src := &ProcCpuinfo{}
	src.Parse(f)
	got, _ := f.Get("phy.cpu.arch_version")
	if got != "x86_64_v0" { // absence yields v0 (no flags matched)
		t.Fatalf("expected x86_64_v0 got %s", got)
	}
}
