"""The main() function implementation for the program."""

import argparse
import logging
import sys

from clu import __about__
from clu.logs import setup_logging
import clu.cmd.report
import clu.cmd.archive
import clu.cmd.requires

log = logging.getLogger(__name__)


def parse_cmdline(args=None):
    parser = argparse.ArgumentParser(prog="clu", description="clu utility")
    parser.add_argument(
        "--version", action="version", version="%%(prog)s %s" % __about__.__version__
    )

    parser.add_argument(
        "--verbose",
        "-v",
        action="count",
        default=0,
        dest="verbosity",
        help="Increase verbosity of output (can be used multiple times).",
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
    setup_logging(args.verbosity)

    log.info("Starting clu utility...")
    log.debug(f"Command line: {args}")
    log.trace(f"tracing is on and working: {args}")

    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
