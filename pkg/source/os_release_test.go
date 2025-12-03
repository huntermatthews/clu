package source_test

import (
	"errors"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/sources"
)

var sampleOsRelease = `NAME="Test Linux"
ID=testlinux
VERSION_ID="1.2"
PRETTY_NAME="Test Linux 1.2"`

// TestOsReleaseProvides validates registration of provided keys.
func TestOsReleaseProvides(t *testing.T) {
	src := &sources.OsRelease{}
	p := pkg.Provides{}
	src.Provides(p)
	for _, k := range []string{"os.distro.name", "os.distro.version"} {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

// TestOsReleaseSuccess simulates a typical os-release file.
func TestOsReleaseSuccess(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return sampleOsRelease, nil }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	src := &sources.OsRelease{}
	src.Parse(f)
	name, _ := f.Get("os.distro.name")
	ver, _ := f.Get("os.distro.version")
	if name != "testlinux" || ver != "1.2" {
		t.Fatalf("unexpected parsed values name=%q ver=%q", name, ver)
	}
}

func TestOsReleaseMissing(t *testing.T) {
	orig := pkg.FileReader
	pkg.FileReader = func(path string) (string, error) { return "", errors.New("missing") }
	defer func() { pkg.FileReader = orig }()
	f := pkg.NewFacts()
	src := &sources.OsRelease{}
	src.Parse(f)
	name, _ := f.Get("os.distro.name")
	ver, _ := f.Get("os.distro.version")
	if name != sources.ParseFailMsg || ver != sources.ParseFailMsg {
		t.Fatalf("expected ParseFailMsg got name=%q ver=%q", name, ver)
	}
}
