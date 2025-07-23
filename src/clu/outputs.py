"""Doc Incomplete."""

import json


from clu.facts import get_fact, get_all_facts


def output_dots():
    for key in sorted(get_all_facts().keys()):
        value = get_fact(key)
        print(f"{key}: {value}")


def output_shell():
    for key in sorted(get_all_facts().keys()):
        value = get_fact(key)
        key_var = key.upper().replace(".", "_")
        print(f"{key_var}=\"{value}\"")


def output_json():
    print(json.dumps(get_all_facts(), indent=2))
