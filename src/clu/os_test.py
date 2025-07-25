
from clu.facts import Facts
from clu.os_linux import parse_os_release, requires_os_release, parse_lscpu, requires_lscpu
from clu.requires import Requires


def requires_os_test(requires: Requires) -> None:
    requires_os_release(requires)
    requires_lscpu(requires)


def parse_os_test(facts: Facts) -> None:
    parse_os_release(facts)
    parse_lscpu(facts)
