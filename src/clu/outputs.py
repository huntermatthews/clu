"""Doc Incomplete."""

import json


from clu import facts
from clu.debug import trace


def output_dots():
    trace("output_dots begin")
    for key in sorted(facts.keys()):
        value = facts[key]
        print(f"{key}: {value}")


def output_shell():
    trace("output_shell begin")
    for key in sorted(facts.keys()):
        value = facts[key]
        key_var = key.upper().replace(".", "_")
        print(f"{key_var}=\"{value}\"")


def output_json():
    trace("output_json begin")

    print(json.dumps(facts, indent=2))
