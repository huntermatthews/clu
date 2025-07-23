"""Doc Incomplete."""

import json


def output_dots(facts):
    for key in sorted(facts.keys()):
        value = facts[key]
        print(f"{key}: {value}")


def output_shell(facts):
    for key in sorted(facts.keys()):
        value = facts[key]
        key_var = key.upper().replace(".", "_")
        print(f"{key_var}=\"{value}\"")


def output_json(facts):
    print(json.dumps(facts, indent=2))
