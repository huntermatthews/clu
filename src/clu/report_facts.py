"""Doc Incomplete."""

from clu import config
from clu.facts import get_fact
from clu.outputs import output_dots, output_shell
from clu.os_darwin import parse_os_darwin
from clu.os_linux import parse_os_linux
from clu.os_test import parse_os_test
from clu.os_unsupported import parse_os_unsupported
from clu.os_generic import parse_uname   # There's always a uname, even if it's mocked.


def do_report_facts():
    """Generate a report based on the current (possibly mocked) OS."""

    if config.test:
        # If we're in test mode, we don't need to do any checks.
        # We just parse the test OS.
        facts = parse_os_test()
    else:
        parse_uname()
        if get_fact("os.kernel.name") == "Darwin":
            facts = parse_os_darwin()
        elif get_fact("os.kernel.name") == "Linux":
            facts = parse_os_linux()
        else:
            facts = parse_os_unsupported()

    if config.output == "dots":
        output_dots(facts)
    elif config.output == "shell":
        output_shell(facts)
    else:
        output_dots(facts)
