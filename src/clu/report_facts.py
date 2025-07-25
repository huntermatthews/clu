"""Doc Incomplete."""

from clu import config
from clu.facts import Facts
from clu.outputs import output_dots, output_shell
from clu.os_darwin import parse_os_darwin
from clu.os_linux import parse_os_linux
from clu.os_test import parse_os_test
from clu.os_unsupported import parse_os_unsupported
from clu.os_generic import parse_uname

def do_report_facts() -> None:
    """Generate a report based on the current (possibly mocked) OS."""

    if config.test:
        # If we're in test mode, we don't need to do any checks.
        # We just parse the test OS.
        facts = Facts()
        parse_os_test(facts)
    else:
        facts = Facts()
        parse_uname(facts)
        if facts["os.kernel.name"] == "Darwin":
            parse_os_darwin(facts)
        elif facts["os.kernel.name"] == "Linux":
            parse_os_linux(facts)
        else:
            parse_os_unsupported(facts)

    if config.output == "dots":
        output_dots(facts)
    elif config.output == "shell":
        output_shell(facts)
    else:
        output_dots(facts)
