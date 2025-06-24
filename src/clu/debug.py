"""Doc Incomplete."""

import sys
import traceback

from clu import config

def panic(*args):
    """
    Print a fatal error message and exit.
    """
    if config.debug > 1:
        traceback.print_stack()

    print("FATAL:", *args, file=sys.stderr)
    sys.exit(1)


def debug(*args):
    """
    Print debug messages if debug state is not 'off'.
    """
    if config.debug > 0:
        print("DEBUG:", *args, file=sys.stderr)


def debug_var(var_name, value):
    """
    Print the value of a variable for debugging.
    """
    if config.debug > 0:
        print(f"DEBUG: variable {var_name} == '{value}'", file=sys.stderr)


def debug_var_list(var_name, var_list):
    """
    Print the contents of a list variable for debugging.
    """
    if config.debug > 0:
        print(
            f"DEBUG: variable `{var_name}` of length {len(var_list)}:", file=sys.stderr
        )
        for i, val in enumerate(var_list, 1):
            print(f"     debug: index: {i}, Value: {val}", file=sys.stderr)


def trace(*args):
    """
    Print trace messages if debug state is 'trace'.
    """
    if config.debug > 1:
        print("TRACE:", *args, file=sys.stderr)


def trace_var_list(var_name, var_list):
    """
    Print the contents of a list variable for tracing.
    """
    if config.debug > 1:
        print(
            f"TRACE: variable `{var_name}` of length {len(var_list)}:", file=sys.stderr
        )
        for i, val in enumerate(var_list, 1):
            print(f"     trace: index: {i}, Value: {val}", file=sys.stderr)

# def caller_function():
#     return (inspect.stack()[1].function, inspect.stack()[1].filename)
