package sources

import (
	"os/exec"
	"strings"
	"testing"

	pkg "github.com/huntermatthews/clu/pkg"
)

func TestUptimeProvides(t *testing.T) {
	u := &Uptime{}
	p := pkg.NewProvides()
	u.Provides(p)
	if _, ok := p["run.uptime"]; !ok {
		t.Errorf("expected run.uptime provide key")
	}
}

func TestUptimeParse(t *testing.T) {
	if _, err := exec.LookPath("uptime"); err != nil {
		t.Skip("uptime not available")
	}
	u := &Uptime{}
	facts := pkg.NewFacts()
	u.Parse(facts)
	if !facts.Contains("run.uptime") {
		t.Errorf("missing run.uptime fact")
	}
}

func TestUptimeRegexPluralSingular(t *testing.T) {
	samples := []struct {
		line string
		want string
	}{
		{"15:42  up 3 days,  2:17, 1 user, load averages: 1.10 1.05 1.00", "3 days,  2:17"},
		{"15:42  up 3 days,  2:17, 2 users, load averages: 1.10 1.05 1.00", "3 days,  2:17"},
	}
	for _, s := range samples {
		m := uptimeRegex.FindStringSubmatch(s.line)
		if len(m) < 2 {
			// Regex failed to match expected pattern
			t.Errorf("no match for line: %q", s.line)
			continue
		}
		got := strings.TrimSuffix(m[1], ",")
		if got != s.want {
			t.Errorf("uptime parse got %q want %q", got, s.want)
		}
	}
}
