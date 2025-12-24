package sources

import (
	"errors"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

func TestProcUptimeProvides(t *testing.T) {
	src := &ProcUptime{}
	p := types.Provides{}
	src.Provides(p)
	if _, ok := p["run.uptime"]; !ok {
		t.Fatalf("missing provides key run.uptime")
	}
}

func TestProcUptimeSuccess(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "12345.67 54321.00", nil }
	defer func() { pkg.FileReader = orig }()
	f := types.NewFacts()
	src := &ProcUptime{}
	src.Parse(f)
	got, _ := f.Get("run.uptime")
	expected := pkg.SecondsToText(12345)
	if got != expected {
		t.Fatalf("want %q got %q", expected, got)
	}
}

func TestProcUptimeMissing(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "", errors.New("missing") }
	defer func() { pkg.FileReader = orig }()
	f := types.NewFacts()
	src := &ProcUptime{}
	src.Parse(f)
	got, _ := f.Get("run.uptime")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}

func TestProcUptimeMalformed(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "not_a_number", nil }
	defer func() { pkg.FileReader = orig }()
	f := types.NewFacts()
	src := &ProcUptime{}
	src.Parse(f)
	got, _ := f.Get("run.uptime")
	if got != types.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg for malformed content got %q", got)
	}
}
