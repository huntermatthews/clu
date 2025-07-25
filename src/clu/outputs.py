"""Doc Incomplete."""

import json

from .facts import Facts

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
