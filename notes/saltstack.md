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

## Sample Version Output

```shell
Salt Version:
          Salt: 3006.19

Python Version:
        Python: 3.10.19 (main, Jan  7 2026, 23:50:47) [GCC 11.2.0]

Dependency Versions:
          cffi: 2.0.0
      cherrypy: 18.10.0
  cryptography: 42.0.5
      dateutil: 2.8.1
     docker-py: Not Installed
         gitdb: Not Installed
     gitpython: Not Installed
        Jinja2: 3.1.6
       libgit2: Not Installed
  looseversion: 1.0.2
      M2Crypto: Not Installed
          Mako: Not Installed
       msgpack: 1.0.2
  msgpack-pure: Not Installed
  mysql-python: Not Installed
     packaging: 24.0
     pycparser: 2.21
      pycrypto: Not Installed
  pycryptodome: 3.19.1
        pygit2: Not Installed
  python-gnupg: 0.4.8
        PyYAML: 6.0.1
         PyZMQ: 23.2.0
        relenv: 0.22.2
         smmap: Not Installed
       timelib: 0.3.0
       Tornado: 4.5.3
           ZMQ: 4.3.4

System Versions:
          dist: ubuntu 24.04.3 noble
        locale: utf-8
       machine: x86_64
       release: 5.15.167.4-microsoft-standard-WSL2
        system: Linux
       version: Ubuntu 24.04.3 noble
```
