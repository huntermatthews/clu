"""Doc Incomplete."""

import json


from clu.facts import get_fact, get_all_facts
from clu.debug import trace


def output_dots():
    trace("output_dots begin")
    for key in sorted(get_all_facts().keys()):
        value = get_fact(key)
        print(f"{key}: {value}")


def output_shell():
    trace("output_shell begin")
    for key in sorted(get_all_facts().keys()):
        value = get_fact(key)
        key_var = key.upper().replace(".", "_")
        print(f"{key_var}=\"{value}\"")


def output_json():
    trace("output_json begin")

    print(json.dumps(get_all_facts(), indent=2))
