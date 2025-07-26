from clu import config
from clu.debug import panic
from clu.facts import Facts
from clu.os_test import requires_os_test, parse_os_test
from clu.os_darwin import requires_os_darwin, parse_os_darwin
from clu.os_linux import requires_os_linux, parse_os_linux
from clu.os_generic import parse_uname


def get_os_functions() -> tuple:
    """Get the requirements and parsing functions for the current OS."""

    if config.test:
        # If we're in test mode, we don't need to do any checks.
        # We just parse the test OS.
        return (requires_os_test, parse_os_test)

    facts = Facts()
    parse_uname(facts)
    kernel_name = facts.get("os.kernel.name", "Unknown")

    if kernel_name == "Darwin":
        return (requires_os_darwin, parse_os_darwin)
    elif kernel_name == "Linux":
        return (requires_os_linux, parse_os_linux)
    else:
        panic("Unsupported OS")
        return (None, None)    # make mypy happy, but this line is unreachable
