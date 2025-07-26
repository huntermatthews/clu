"""Doc Incomplete."""

from clu import config
from clu.outputs import output_dots, output_shell
from clu.os_map import get_os_functions


def do_report_facts() -> None:
    """Generate a report based on the current OS."""

    (_, parse_fn) = get_os_functions()
    print(type(parse_fn))
    facts = parse_fn()
    print(type(facts))

    if config.output == "dots":
        output_dots(facts)
    elif config.output == "shell":
        output_shell(facts)
    else:
        output_dots(facts)
