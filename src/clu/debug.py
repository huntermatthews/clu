"""Doc Incomplete."""

import sys
import traceback

_debug = "off"


def panic(*args):
    """
    Print a fatal error message and exit.
    """
    if _debug == "trace":
        traceback.print_stack()

    print("FATAL:", *args, file=sys.stderr)
    sys.exit(1)


def debug_state(state=None):
    """
    Set or query the debug state.
    if state is passed, it can be 'off', 'debug', 'trace'.
    the current state is returned regardless.
    """

    if state in ("off", "debug", "trace"):
        global _debug
        _debug = state
    elif state is not None:
        panic(
            f"Invalid debug state '{state}'. Valid states are 'off', 'debug', 'trace'."
        )
    return _debug


def debug(*args):
    """
    Print debug messages if debug state is not 'off'.
    """
    if _debug != "off":
        print("DEBUG:", *args, file=sys.stderr)


def debug_var(var_name, value):
    """
    Print the value of a variable for debugging.
    """
    if _debug != "off":
        print(f"DEBUG: variable {var_name} == '{value}'", file=sys.stderr)


def debug_var_list(var_name, var_list):
    """
    Print the contents of a list variable for debugging.
    """
    if _debug != "off":
        print(
            f"DEBUG: variable `{var_name}` of length {len(var_list)}:", file=sys.stderr
        )
        for i, val in enumerate(var_list, 1):
            print(f"     debug: index: {i}, Value: {val}", file=sys.stderr)


def trace(*args):
    """
    Print trace messages if debug state is 'trace'.
    """
    if _debug == "trace":
        print("TRACE:", *args, file=sys.stderr)


def trace_var_list(var_name, var_list):
    """
    Print the contents of a list variable for tracing.
    """
    if _debug == "trace":
        print(
            f"TRACE: variable `{var_name}` of length {len(var_list)}:", file=sys.stderr
        )
        for i, val in enumerate(var_list, 1):
            print(f"     trace: index: {i}, Value: {val}", file=sys.stderr)
