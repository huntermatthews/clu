// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/huntermatthews/clu/pkg/facts"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/global"
	"github.com/huntermatthews/clu/pkg/input"
)

const (
	maxNoSaltAge     = 60 * 24 * time.Hour // 2 months
	maxHighstateAge  = 2 * 24 * time.Hour  // 2 days
	maxKnownHostsAge = 30 * 24 * time.Hour // 30 days
)

type CheckCmd struct{}

func (f *CheckCmd) Run(stdout input.Stdout, stderr input.Stderr) error {
	slog.Info("Running CheckCmd...")

	if global.Config.MockDir != "" {
		return fmt.Errorf("check subcommand does not yet support mock mode")
	}

	osys := facts.OpSysFactory()
	provides := osys.Provides()
	facts := types.NewFactDB()

	factNames := []string{
		"phy.platform",
		"sys.asset_no",
		"salt.no_salt.exists",
		"salt.no_salt.age",
		"salt.no_salt.reason",
		"salt.highstate.age",
	}

	// parse facts required for checks
	parseFactsBySpecs(provides, facts, factNames)

	// Check: physical host and no asset tag
	if err := checkPlatformAssetNo(facts); err != nil {
		return err
	}

	// Check: Salt minion checks
	if err := checkSaltMinion(facts); err != nil {
		return err
	}

	// Check: zabbix client running
	if err := checkZabbixClient(); err != nil {
		return err
	}

	// Check: /etc/ssh/ssh_known_hosts too old
	if err := checkSSHKnownHosts(); err != nil {
		return err
	}

	fmt.Fprintln(stdout, "All checks passed")
	return nil

}

func checkSSHKnownHosts() error {
	knownHostsTime, err := GetSSHKnownHostsAge()
	if err != nil {
		// File doesn't exist or can't be read - silently pass
		return nil
	}
	// Check if file is older than maxKnownHostsAge
	age := time.Since(knownHostsTime)
	if age > maxKnownHostsAge {
		return fmt.Errorf("/etc/ssh/ssh_known_hosts is too old (%.0f days)", age.Hours()/24)
	}
	return nil
}

func checkSaltMinion(facts *types.FactDB) error {
	// Check: /no_salt exists
	noSaltExists, noSaltOK := facts.Get("salt.no_salt.exists")
	if !noSaltOK {
		panic("required fact salt.no_salt.exists not found")
	}
	if noSaltExists == "True" {
		slog.Warn("/no_salt file exists - salt may be disabled on this host")
	}

	// Check: /no_salt is too old
	if noSaltExists == "True" {
		noSaltAge, ageOK := facts.Get("salt.no_salt.age")
		if !ageOK {
			panic("required fact salt.no_salt.age not found")
		}
		// Parse the RFC3339 timestamp with space separator
		fileTime, err := time.Parse("2006-01-02 15:04:05Z07:00", noSaltAge)
		if err != nil {
			slog.Error("Failed to parse /no_salt age", "error", err)
			panic("Failed to parse /no_salt age")
		}
		// Check if file is older than maxNoSaltAge
		age := time.Since(fileTime)
		if age > maxNoSaltAge {
			return fmt.Errorf("/no_salt file is too old (%.0f days)", age.Hours()/24)
		}
	}

	// Check: /no_salt exists but no message
	if noSaltExists == "True" {
		noSaltReason, reasonOK := facts.Get("salt.no_salt.reason")
		if !reasonOK {
			panic("required fact salt.no_salt.reason not found")
		}
		if noSaltReason == "" {
			return fmt.Errorf("/no_salt exists but has no message")
		}
	}

	// Check: highstate too old
	highstateAge, highstateOK := facts.Get("salt.highstate.age")
	if !highstateOK {
		panic("required fact salt.highstate.age not found")
	}

	if highstateAge != "" {
		// Parse the RFC3339 timestamp with space separator
		fileTime, err := time.Parse("2006-01-02 15:04:05Z07:00", highstateAge)
		if err != nil {
			slog.Error("Failed to parse highstate age", "error", err)
			panic("Failed to parse highstate age")
		}
		// Check if file is older than maxHighstateAge
		age := time.Since(fileTime)
		if age > maxHighstateAge {
			return fmt.Errorf("Highstate file is too old (%.0f days)", age.Hours()/24)
		}
	}

	return nil
}

func checkPlatformAssetNo(facts *types.FactDB) error {
	platformFact, platformOK := facts.Get("phy.platform")
	if !platformOK {
		panic("required fact phy.platform not found")
	}
	assetNoFact, assetOK := facts.Get("sys.asset_no")
	if !assetOK {
		panic("required fact sys.asset_no not found")
	}
	if platformFact == "physical" && (assetNoFact == "unknown" || assetNoFact == "") {
		return fmt.Errorf("Host is physical and has no asset tag")
	}
	return nil
}

// IsServiceEnabled checks if a systemd service is enabled.
// Returns true if enabled, false otherwise, and any error encountered.
func IsServiceEnabled(serviceName string) (bool, error) {
	output, rc, err := input.CommandRunner("systemctl is-enabled " + serviceName)
	if err != nil {
		return false, fmt.Errorf("failed to check if %s is enabled: %w", serviceName, err)
	}

	if rc == 0 && strings.TrimSpace(output) == "enabled" {
		return true, nil
	}

	return false, nil
}

// IsServiceActive checks if a systemd service is active (started and running).
// Returns true if active, false otherwise, and any error encountered.
func IsServiceActive(serviceName string) (bool, error) {
	output, rc, err := input.CommandRunner("systemctl is-active " + serviceName)
	if err != nil {
		return false, fmt.Errorf("failed to check if %s is active: %w", serviceName, err)
	}

	if rc == 0 && strings.TrimSpace(output) == "active" {
		return true, nil
	}

	return false, nil
}

// checkZabbixClient checks if the zabbix-agent service is both enabled and active.
// Returns an error if the service is not enabled or not active.
func checkZabbixClient() error {
	enabled, err := IsServiceEnabled("zabbix-agent")
	if err != nil {
		return fmt.Errorf("error checking if zabbix-agent is enabled: %w", err)
	}
	if !enabled {
		return fmt.Errorf("zabbix agent is not enabled")
	}

	active, err := IsServiceActive("zabbix-agent")
	if err != nil {
		return fmt.Errorf("error checking if zabbix-agent is active: %w", err)
	}
	if !active {
		return fmt.Errorf("zabbix agent is not active")
	}

	return nil
}

// GetSSHKnownHostsAge retrieves the modification time of /etc/ssh/ssh_known_hosts.
// Returns the time.Time and any error encountered.
func GetSSHKnownHostsAge() (time.Time, error) {
	return input.FileAgeReader("/etc/ssh/ssh_known_hosts")
}

/*
CODE ANALYSIS - Issues to Address Later:

1. Logic Issue: Redundant /no_salt conditionals (lines 96-122)
   - All three checks for /no_salt age and message are guarded by if noSaltExists == "True"
   - The facts salt.no_salt.age and salt.no_salt.reason are always retrieved via parseFactsBySpecs()
     regardless of file existence
   - If /no_salt doesn't exist, you fetch and panic-check for facts that likely don't exist
   - Recommendation: Either remove these facts from the factNames list when the file doesn't exist,
     or handle missing facts gracefully instead of panicking

2. Inconsistent Error Handling Philosophy
   - Missing facts → panic (developer/setup error)
   - Parse errors → log + panic (data format error)
   - Check failures → return error (operational issue)
   - This is actually fairly consistent, but it means the factNames list should be dynamic based
     on what files exist

3. Empty highstate age handling (lines 125-141)
   - Line 127 checks if highstateAge != "" before parsing
   - But line 124 already panics if the fact doesn't exist
   - This means an empty string would pass silently - is that intentional?
   - Question: Does empty highstateAge mean the file doesn't exist, or is it a valid state?

4. checkSSHKnownHosts silently passes on any error (lines 70-72)
   - Comment says "File doesn't exist or can't be read - silently pass"
   - This could hide permission errors, filesystem issues, etc.
   - Consider: Only silently pass on os.IsNotExist(), panic or log other errors

5. Error message capitalization inconsistency
   - Line 152: "Host is physical and has no asset tag" (capital H)
   - Line 109/117/138: lowercase starts
   - Mostly fixed from earlier, but line 152 still has capital start

6. Constants comment mismatch (line 16)
   - maxNoSaltAge = 60 * 24 * time.Hour // 2 months
   - 60 days ≠ 2 months (months vary between 28-31 days)
   - Recommendation: Change comment to "60 days" or adjust constant to approximate months

7. checkSaltMinion could accumulate multiple errors
   - Currently returns on first error found
   - If both age AND message are problems, only the first is reported
   - Consider: Accumulate errors and return them all, or document that checks are short-circuiting

8. Missing stderr usage
   - Run() accepts stderr input.Stderr but never uses it
   - All errors go to return value (printed by Kong)
   - slog messages go to configured output
   - Consider: Remove unused parameter or document why it's there
*/
