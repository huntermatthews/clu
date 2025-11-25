"""The main() function implementation for the program."""

import argparse
import logging
import sys

from clu import __about__
from clu.config import set_config, get_config
import clu.cmd.report
import clu.cmd.archive
import clu.cmd.requires

log = logging.getLogger(__name__)
cfg = get_config()


def setup_logging(debug_level) -> None:
    DEBUG_MAP = {
        0: logging.WARNING,  # This is the default debug_level.
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


def parse_cmdline() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        prog="clu",
        description="clu utility",
        #        epilog="end tezt here",
        usage="%(prog)s [options] <command> [<commmand args>]",
    )
    parser.add_argument(
        "--version", action="version", version="%%(prog)s %s" % clu.__about__.__version__
    )
    parser._optionals.title = (
        "options"  # yeah, it's a private member - should have a public way to do this
    )

    parser.add_argument(
        "--debug",
        action="count",
        default=0,
        dest="debug_level",
        help="Increase debugging output (can be used twice).",
    )

    parser.add_argument(
        "--verbose",
        action="count",
        default=1,
        dest="verbose_level",
        help="Increase verbosity of output.",
    )

    subparsers = parser.add_subparsers(
        dest="cmd",
        description="Clu's functionality is dividedi into commands. Use 'clu <command> --help' or the man page for more info.",
        metavar="",  # setting this to empty string removes the ugly {foo,bar} from the help output
        title="commands",
    )
    parser.set_defaults(cmd="report", func=clu.cmd.report.report_facts)

    clu.cmd.archive.parse_args(subparsers)
    clu.cmd.report.parse_args(subparsers)
    clu.cmd.requires.parse_args(subparsers)

    return parser.parse_args()


def main() -> int:
    if sys.version_info < clu.__about__.__minimum_python_info__:
        # Can't use logging here
        print(f"ERROR: Must use at least python {__about__.__minimum_python__}", file=sys.stderr)
        sys.exit(1)

    args = parse_cmdline()
    set_config(args)
    setup_logging(args.debug_level)

    log.info("Starting clu utility...")
    log.debug(f"Command line: {args}")

    return args.func()


if __name__ == "__main__":
    sys.exit(main())
