"""Doc Incomplete."""

import logging

from clu import Provides
from clu import Requires
from clu import Facts
from clu.debug import panic
from clu.input import text_program
from clu.os_generic import (
    provides_uname,
    requires_uname,
    parse_uname,
    provides_uptime,
    requires_uptime,
    parse_uptime,
    provides_clu,
    requires_clu,
    parse_clu,
)

log = logging.getLogger(__name__)


def default_facts_os_darwin() -> list:
    return ["os.name", "os.hostname", "os.version", "os.codename", "run.uptime", "clu.version"]


def provides_os_darwin() -> Provides:
    """Define the provider map for macOS (Darwin)."""
    provides = Provides()

    provides_uname(provides)
    provides_sw_vers(provides)
    provides_macos_name(provides)
    provides_uptime(provides)
    provides_clu(provides)

    return provides


def requires_os_darwin() -> Requires:
    """Define the requirements for macOS (Darwin)."""
    requires = Requires()

    requires_uname(requires)
    requires_sw_vers(requires)
    requires_macos_name(requires)
    requires_uptime(requires)
    requires_clu(requires)
    requires_systemversion_plist(requires)

    return requires


def parse_os_darwin() -> Facts:
    """Parse the facts for macOS (Darwin)."""
    facts = Facts()

    # Nothing explicitly says Apple, but we know its apple because Darwin is the OS
    facts["sys.vendor"] = "Apple"

    parse_uname(facts)
    parse_sw_vers(facts)
    parse_macos_name(facts)
    parse_uptime(facts)
    parse_clu(facts)

    return facts
