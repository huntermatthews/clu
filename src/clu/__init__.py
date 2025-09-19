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
