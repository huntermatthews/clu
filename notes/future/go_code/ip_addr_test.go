package sources

import (
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
	facts "github.com/huntermatthews/clu/pkg/facts"
)

func TestIpAddrProvides(t *testing.T) {
	src := &IpAddr{}
	p := facts.Provides{}
	src.Provides(p)
	for _, k := range []string{"net.macs", "net.ipv4", "net.ipv6", "net.devs"} {
		if _, ok := p[k]; !ok {
			t.Fatalf("missing provides key %s", k)
		}
	}
}

func TestIpAddrSuccess(t *testing.T) {
	json := `[
        {"ifname":"eth0","address":"00:11:22:33:44:55","addr_info":[{"family":"inet","local":"192.168.1.10"},{"family":"inet6","local":"fe80::1234"}]},
        {"ifname":"lo","address":"","addr_info":[{"family":"inet","local":"127.0.0.1"}]}
    ]`
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return json, 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := facts.NewFacts()
	src := &IpAddr{}
	src.Parse(f)
	cases := map[string]string{
		"net.devs": "eth0 lo ",
		"net.macs": "00:11:22:33:44:55 ",
		"net.ipv4": "192.168.1.10 127.0.0.1 ",
		"net.ipv6": "fe80::1234 ",
	}
	for k, want := range cases {
		got, _ := f.Get(k)
		if got != want {
			t.Fatalf("%s want %q got %q", k, want, got)
		}
	}
}

func TestIpAddrFailureRC(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "", 1 }
	defer func() { pkg.CommandRunner = orig }()
	f := facts.NewFacts()
	src := &IpAddr{}
	src.Parse(f)
	for _, k := range []string{"net.macs", "net.ipv4", "net.ipv6", "net.devs"} {
		got, _ := f.Get(k)
		if got != ParseFailMsg {
			t.Fatalf("%s expected ParseFailMsg got %q", k, got)
		}
	}
}

func TestIpAddrEmptyOutput(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := facts.NewFacts()
	src := &IpAddr{}
	src.Parse(f)
	for _, k := range []string{"net.macs", "net.ipv4", "net.ipv6", "net.devs"} {
		got, _ := f.Get(k)
		if got != ParseFailMsg {
			t.Fatalf("%s expected ParseFailMsg got %q", k, got)
		}
	}
}

func TestIpAddrMalformedJSON(t *testing.T) {
	orig := pkg.CommandRunner
	pkg.CommandRunner = func(cmd string) (string, int) { return "not json", 0 }
	defer func() { pkg.CommandRunner = orig }()
	f := facts.NewFacts()
	src := &IpAddr{}
	src.Parse(f)
	for _, k := range []string{"net.macs", "net.ipv4", "net.ipv6", "net.devs"} {
		got, _ := f.Get(k)
		if got != ParseFailMsg {
			t.Fatalf("%s expected ParseFailMsg got %q", k, got)
		}
	}
}
