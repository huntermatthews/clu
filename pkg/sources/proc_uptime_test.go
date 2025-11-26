package sources_test

import (
	"errors"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

func TestProcUptimeProvides(t *testing.T) {
	src := &sources.ProcUptime{}
	p := pkg.Provides{}
	src.Provides(p)
	if _, ok := p["run.uptime"]; !ok {
		t.Fatalf("missing provides key run.uptime")
	}
}

func TestProcUptimeSuccess(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "12345.67 54321.00", nil }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	src := &sources.ProcUptime{}
	src.Parse(f)
	got, _ := f.Get("run.uptime")
	// 12345 seconds -> expected SecondsToText output
	expected := pkg.SecondsToText(12345)
	if got != expected {
		t.Fatalf("want %q got %q", expected, got)
	}
}

func TestProcUptimeMissing(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "", errors.New("missing") }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	src := &sources.ProcUptime{}
	src.Parse(f)
	got, _ := f.Get("run.uptime")
	if got != sources.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got %q", got)
	}
}

func TestProcUptimeMalformed(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "not_a_number", nil }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	src := &sources.ProcUptime{}
	src.Parse(f)
	got, _ := f.Get("run.uptime")
	if got != sources.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg for malformed content got %q", got)
	}
}
