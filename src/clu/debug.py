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

def debug_var(var_name, var_value):
    if config.debug > 0:
        _debug_var("DEBUG", var_name, var_value)

def trace(*args):
    """
    Print trace messages if debug state is 'trace'.
    """
    if config.debug > 1:
        print("TRACE:", *args, file=sys.stderr)

def trace_var(var_name, var_value):
    if config.debug > 1:
        _debug_var("TRACE", var_name, var_value)

def _debug_var(level, var_name, var_value):
    """
    Print the value of a variable for debugging.
    """
    if config.debug > 0:
        var_type = type(var_value)
        if var_type is str:
            print(f"{level}: {var_name} (str[{len(var_value)}]) == '{var_value}'", file=sys.stderr)
        elif var_type is list or var_type is tuple:
            print(f"{level}: {var_type}[{len(var_value)}] {var_name}:", file=sys.stderr)
            for i, val in enumerate(var_value, 1):
                print(f"     debug: index: {i}, Value: {val}", file=sys.stderr)
        elif var_type is dict:
            print(f"{level}: dict[{len(var_value)}] {var_name}:", file=sys.stderr)
            for key in var_value:
                print(f"  {var_name}[{key}] = {var_value[key]}", file=sys.stderr)
        elif var_type is None:
            print(f"{level}: {var_name} is None...", file=sys.stderr)
        else:
            print(f"{level}: {var_name} (type {var_type} ) == '{var_value}'", file=sys.stderr)


# def caller_function():
#     return (inspect.stack()[1].function, inspect.stack()[1].filename)
