
from clu.facts import Facts
from clu.os_linux import parse_uname, requires_os_release, parse_lscpu, requires_lscpu
from clu.provides import Provides
from clu.requires import Requires




def provides_os_test() -> Provides:
    provides = Provides()

    return provides

def requires_os_test() -> Requires:
    """Define the requirements for the test OS."""
    requires = Requires()

    requires_os_release(requires)
    requires_lscpu(requires)
    return requires


def parse_os_test() -> Facts:
    """Parse the facts for the test OS."""
    facts = Facts()

    parse_uname(facts)
    parse_lscpu(facts)
    return facts
