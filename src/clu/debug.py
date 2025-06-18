import sys

# Global debug state: 'off', 'debug', or 'trace'
_debug = "off"


def debug_state(cmd):
    """
    Set or query the debug state.
    Usage: debug_state('off'|'debug'|'trace'|'status'|'print_status')
    """
    global _debug
    if cmd in ("off", "debug", "trace"):
        _debug = cmd
    elif cmd == "status":
        print(_debug)
    elif cmd == "print_status":
        print(f"DEBUG: status is {_debug}")
    else:
        print(f"ERROR: Unknown arg {cmd}", file=sys.stderr)
        print("Usage: debug_state off|debug|trace|status|print_status", file=sys.stderr)


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
    if var_name == "argv":
        print("ERROR: You can't debug_var argv itself...", file=sys.stderr)
        sys.exit(11)
    if _debug != "off":
        print(f"DEBUG: variable {var_name} == '{value}'", file=sys.stderr)


def debug_var_list(var_name, var_list):
    """
    Print the contents of a list variable for debugging.
    """
    if _debug != "off":
        count = len(var_list)
        if count == 0:
            print(
                f"DEBUG: variable `{var_name}` of length {count} == '{var_list}'",
                file=sys.stderr,
            )
        else:
            print(f"DEBUG: variable `{var_name}` of length {count}:", file=sys.stderr)
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
        count = len(var_list)
        if count == 0:
            print(
                f"TRACE: variable `{var_name}` of length {count} == '{var_list}'",
                file=sys.stderr,
            )
        else:
            print(f"TRACE: variable `{var_name}` of length {count}:", file=sys.stderr)
            for i, val in enumerate(var_list, 1):
                print(f"     trace: index: {i}, Value: {val}", file=sys.stderr)
