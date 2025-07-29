"""Doc Incomplete."""

import json

from clu import config
from clu.os_map import get_os_functions
from clu.facts import Facts


def do_report_facts() -> None:
    """Generate a report based on the current OS."""

    (_, parse_fn) = get_os_functions()
    facts = parse_fn()

    if config.output == "json":
        output_json(facts)
    elif config.output == "shell":
        output_shell(facts)
    else:
        output_dots(facts)


def output_dots(facts: Facts) -> None:
    for key in facts:
        value = facts[key]
        print(f"{key}: {value}")


def output_shell(facts: Facts) -> None:
    for key in facts:
        value = facts[key]
        key_var = key.upper().replace(".", "_")
        print(f"{key_var}=\"{value}\"")


def output_json(facts: Facts) -> None:
    print(json.dumps(facts, indent=2))
