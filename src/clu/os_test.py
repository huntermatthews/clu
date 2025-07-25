
from clu.facts import Facts
from clu.os_linux import parse_os_release, requires_os_release, parse_lscpu, requires_lscpu
from clu.requires import Requires


def requires_os_test() -> Requires:
    """Define the requirements for the test OS."""
    requires = Requires()

    requires_os_release(requires)
    requires_lscpu(requires)
    return requires


def parse_os_test() -> Facts:
    """Parse the facts for the test OS."""
    facts = Facts()

    parse_os_release(facts)
    parse_lscpu(facts)
    return facts
