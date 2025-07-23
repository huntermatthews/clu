
from clu.os_linux import parse_os_release, requires_os_release


def requires_os_test():
    requires = requires = {
        "files": [],
        "programs": [],
        "apis": [],
    }
    requires.update(requires_os_release())
    return requires


def parse_os_test():
    facts = {}
    facts.update(parse_os_release())
    return facts
