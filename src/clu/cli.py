"""The main() function implementation for the program."""

import argparse
import logging
import sys

from clu import __about__

import clu.cmd.report
import clu.cmd.archive
import clu.cmd.requires


log = logging.getLogger(__name__)


def setup_logging(debug_level) -> None:
    DEBUG_MAP = {
        0: logging.WARNING,  # This is the default debug_level, IE none.
        1: logging.INFO,
        2: logging.DEBUG,
    }

    """Set up the logging configuration based on command line arguments."""

    if debug_level in DEBUG_MAP:
        level = DEBUG_MAP[debug_level]
    else:
        print("ERROR: verbosity levels wrong, using 2", file=sys.stderr)
        level = DEBUG_MAP[2]

    logging.basicConfig(level=level)


def parse_cmdline(args=None):
    parser = argparse.ArgumentParser(prog="clu", description="clu utility")
    parser.add_argument(
        "--version", action="version", version="%%(prog)s %s" % __about__.__version__
    )

    parser.add_argument(
        "--debug",
        action="count",
        default=0,
        dest="debug_level",
        help="Increase debugging output (can be used twice).",
    )
    subparsers = parser.add_subparsers(dest="cmd")
    parser.set_defaults(cmd="report", func=clu.cmd.report.report_facts)

    clu.cmd.archive.parse_args(subparsers)
    clu.cmd.report.parse_args(subparsers)
    clu.cmd.requires.parse_args(subparsers)

    return parser.parse_args()


def main() -> int:
    if sys.version_info < __about__.__minimum_python_info__:
        # Can't use logging here
        print(f"ERROR: Must use at least python {__about__.__minimum_python__}", file=sys.stderr)
        sys.exit(1)

    args = parse_cmdline()
    setup_logging(args.debug_level)

    log.info("Starting clu utility...")
    log.debug(f"Command line: {args}")

    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
