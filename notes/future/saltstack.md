# Saltstack

## Source of Truth

`formulas/zabbix/files/external_scripts/check_saltstack.sh`
This is what zabbix trusts, this is what clu should trust.

Summary: `/var/run/salt/minion/latest_salt_highstate` should be less than 24 hours old
    created by the saltcall-runner.
