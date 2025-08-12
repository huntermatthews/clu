"""The main() function implementation for the program."""

import argparse
import logging
import sys

# import clu
from clu import __about__
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
        "--debug",
        action="store_const",
        const=logging.DEBUG,
        default=logging.WARNING,
        dest="verbosity",
        help="Debug mode -- boring nerdy details.",
    )
    parser.add_argument(
        "--verbose",
        action="store_const",
        const=logging.INFO,
        dest="verbosity",
        help="Verbose mode -- slightly more information than a normal run.",
    )

    subparsers = parser.add_subparsers(dest="cmd")
    parser.set_defaults(cmd="report", func=clu.cmd.report.report_facts)

    clu.cmd.archive.parse_args(subparsers)
    clu.cmd.report.parse_args(subparsers)
    clu.cmd.requires.parse_args(subparsers)

    return parser.parse_args()


def main() -> int:
    if sys.version_info < __about__.__minimum_python_info__:
        print(f"Must use at least python {__about__.__minimum_python__}", file=sys.stderr)
        sys.exit(1)

    args = parse_cmdline()

    if args.verbosity is not logging.WARNING:
        logging.basicargs(level=args.verbosity)

    log.info(f"Starting clu utility... {sys.argv=}")
    log.debug(f"Command line: {args}")

    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
