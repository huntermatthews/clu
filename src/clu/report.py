"""Doc Incomplete."""

from clu import facts, config
from clu.outputs import output_dots, output_shell
from clu.parse_darwin import parse_os_darwin
from clu.parse_linux import parse_os_linux
from clu.parse_generic import parse_os_test, parse_os_unsupported, parse_uname


def do_report():
    """Generate a report based on the current (possibly mocked) OS."""

    parse_uname()
    if facts["os.kernel.name"] == "Darwin":
        parse_os_darwin()
    elif facts["os.kernel.name"] == "Linux":
        # parse_os_linux()
        parse_os_test()
    else:
        parse_os_unsupported()

    if config.output == "dots":
        output_dots()
    elif config.output == "shell":
        output_shell()
    else:
        output_dots()
