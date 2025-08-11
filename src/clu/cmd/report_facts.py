"""Doc Incomplete."""

import json

from clu import config
from clu.os_map import get_os_functions
from clu import Facts


def do_report_facts() -> None:
    """Generate a report based on the current OS."""

    (_, parse_fn, provides_fn, default_facts_fn) = get_os_functions()

    provides_map = provides_fn()
    parsers_to_call = set()

    if config.all:
        config.facts = "os sys phy run salt clu"
    elif not config.facts:
        config.facts = default_facts_fn()

    # Loop through the facts that were requested on the command line and get a set of parsers that will
    # obtain those facts (there's likely duplicates, so we use a set here)
    for fact_spec in config.facts:
        for key in provides_map:
            if key.startswith(fact_spec):
                parsers_to_call.add(provides_map[key])

    # Call the parsers that we found in the previous loop
    parsed_facts = Facts()
    for parser_fn in parsers_to_call:
        parser_fn(parsed_facts)

    # Loop through the facts that were requested on the command line (again) and make a new
    # facts list/dict that has JUST those matching facts
    output_facts = Facts()
    for fact_spec in config.facts:
        for key in parsed_facts:
            if key.startswith(fact_spec):
                output_facts[key] = parsed_facts[key]

    if config.output == "json":
        output_json(output_facts)
    elif config.output == "shell":
        output_shell(output_facts)
    else:
        output_dots(output_facts)


def output_dots(facts: Facts) -> None:
    for key in facts:
        value = facts[key]
        print(f"{key}: {value}")


def output_shell(facts: Facts) -> None:
    for key in facts:
        value = facts[key]
        key_var = key.upper().replace(".", "_")
        print(f'{key_var}="{value}"')


def output_json(facts: Facts) -> None:
    print(json.dumps(facts, indent=2))
