def trace(msg):
    print(f"TRACE: {msg}")


def output_dots():
    trace("output_dots begin")
    for key in sorted(ATTRS.keys()):
        value = ATTRS[key]
        print(f"{key}: {value}")


def output_shell():
    trace("output_shell begin")
    for key in sorted(ATTRS.keys()):
        value = ATTRS[key]
        key_var = key.upper().replace(".", "_")
        print(f"{key_var}='{value}'")


def output_json():
    trace("output_json begin")
    import json

    print(json.dumps(ATTRS, indent=2))


# Example ATTRS dictionary for testing
# ATTRS = {
#     "sys.vendor": "Apple",
#     "cpu.count": 8,
#     "mem.total": "16GB"
