import logging
import sys


def panic(*args):
    """
    Print a fatal error message and exit.
    """

    log = logging.getLogger(__name__)
    log.debug("stack trace:", exc_info=True, stack_info=True)
    log.fatal(" ".join(str(arg) for arg in args))
    sys.exit(1)


def _debug_var(var_name, var_value):
    """
    Return a debug string representation of a variable.
    """

    var_type = type(var_value)
    if var_type is str:
        return f"{var_name} (str[{len(var_value)}]) == '{var_value}'"
    elif var_type is list or var_type is tuple:
        lines = [f"{var_type}[{len(var_value)}] {var_name}:"]
        for i, val in enumerate(var_value, 1):
            lines.append(f"     debug: index: {i}, Value: {val}")
        return "\n".join(lines)
    elif var_type is dict:
        lines = [f"dict[{len(var_value)}] {var_name}:"]
        for key in var_value:
            lines.append(f"  {var_name}[{key}] = {var_value[key]}")
        return "\n".join(lines)
    elif var_type is None:
        return f"{var_name} is None..."
    else:
        return f"{var_name} (type {var_type} ) == '{var_value}'"


# def caller_function():
#     return (inspect.stack()[1].function, inspect.stack()[1].filename)
