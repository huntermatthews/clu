
# hostnames

Based on <https://itbwiki.nhgri.nih.gov/wiki/index.php/Workstation_Naming_Conventions>
^HG-(0\d{7})-[DLTVC][WML]\d$  with gmi flags (global, multiline [^ and $ work], case insensitive)
is the regex for our hostnames.


TODO: make the trailing digit optional - they changed the spec in 2021-12

HG_HOSTNAME_RE = re.compile('^HG-(0\d{7})-[DLTVC][WML](\d)?$', re.ASCII | re.IGNORECASE | re.MULTILINE)


## macOS

nvram -- fmm-computer-name should be there and parseable on NHGRI equipment.


        # FIX: why does smc have multiple id numbers?
        # DOC: "sim" is "Supermicro intelligent management" -- I love the irony
        bmc_type = {
            '674': 'drac',
            '11': 'ilo',
            '10876': 'sim',
            '47196': 'sim',
            '47488': 'sim',
        }

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
