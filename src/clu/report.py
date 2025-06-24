"""Doc Incomplete."""

from clu import facts, config
from clu.outputs import output_dots, output_shell
from clu.os_darwin import parse_os_darwin
from clu.os_linux import parse_os_linux
from clu.os_test import parse_os_test
from clu.os_unsupported import parse_os_unsupported
from clu.os_generic import parse_uname   # There's always a uname, even if it's mocked.

def do_report():
    """Generate a report based on the current (possibly mocked) OS."""

    if config.test:
        # If we're in test mode, we don't need to do any checks.
        # We just parse the test OS.
        parse_os_test()
    else:
        parse_uname()
        if facts["os.kernel.name"] == "Darwin":
            parse_os_darwin()
        elif facts["os.kernel.name"] == "Linux":
            parse_os_linux()
        else:
            parse_os_unsupported()

    if config.output == "dots":
        output_dots()
    elif config.output == "shell":
        output_shell()
    else:
        output_dots()
