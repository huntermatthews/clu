from clu import config, panic
from clu.facts import Facts
from clu.os_test import requires_os_test, parse_os_test, provides_os_test, default_facts_os_test
from clu.os_darwin import (
    default_facts_os_darwin,
    requires_os_darwin,
    parse_os_darwin,
    provides_os_darwin,
)
from clu.os_linux import (
    default_facts_os_linux,
    requires_os_linux,
    parse_os_linux,
    provides_os_linux,
)
from clu.os_generic import parse_uname


def get_os_functions() -> tuple:
    """Get the requirements and parsing functions for the current OS."""

    if config.test:
        # If we're in test mode, we don't need to do any checks.
        # We just use the test OS.
        return (requires_os_test, parse_os_test, provides_os_test, default_facts_os_test)

    facts = Facts()
    parse_uname(
        facts
    )  # we end up throwing this away, but thats temporary until full pytest support
    kernel_name = facts.get("os.kernel.name", "Unknown")

    if kernel_name == "Darwin":
        return (requires_os_darwin, parse_os_darwin, provides_os_darwin, default_facts_os_darwin)
    elif kernel_name == "Linux":
        return (requires_os_linux, parse_os_linux, provides_os_linux, default_facts_os_linux)
    else:
        panic("Unsupported OS")
        return (None, None)  # make mypy happy, but this line is unreachable
