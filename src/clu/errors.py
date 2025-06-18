import sys
import traceback


def panic(*args):
    """
    Print a fatal error message and exit.
    If debug_state is 'trace', print a stack trace.
    """
    # You should implement debug_state() elsewhere to return 'off', 'debug', or 'trace'
    try:
        from debug import debug_state

        debug = debug_state("status")
    except ImportError:
        debug = "off"

    if debug == "trace":
        traceback.print_stack()

    print("FATAL:", *args, file=sys.stderr)
    sys.exit(1)


def command_not_found(cmd):
    """
    Handler for command-not-found situations.
    """
    print(f"FATAL: Command not found: {cmd}", file=sys.stderr)
    sys.exit(12)
