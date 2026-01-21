# Saltstack

## Source of Truth

`formulas/zabbix/files/external_scripts/check_saltstack.sh`
This is what zabbix trusts, this is what clu should trust.

Summary: `/var/run/salt/minion/latest_salt_highstate` should be less than 24 hours old
    created by the saltcall-runner.


## is host salted at all

```python
def is_salt_host():
    """ Is this a salt host - master, syndic or minion?
    """
    salt_dir = Path('/etc/salt')
    return salt_dir.exists()


        self.reg.insert(Trait(self.group_name + 'daemons.machine_id', self._get_machine_id))
        self.reg.insert(Trait(self.group_name + 'daemons.minion_id', self._get_minion_id))

    def _get_machine_id(self, _):
        return read_simple_file('/etc/machine-id')

    def _get_minion_id(self, _):
        return read_simple_file('/etc/salt/minion_id')
```
