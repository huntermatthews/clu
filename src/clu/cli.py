"""The main() function implementation for the program."""

import argparse
import logging
import sys

# import clu
from clu import __about__, config
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

    subparsers = parser.add_subparsers(dest="cmd") #, required=True)
    parser.set_defaults(cmd="report", func=clu.cmd.report.run)   # This WORKS!

    clu.cmd.archive.setup_args(subparsers)
    clu.cmd.report.setup_args(subparsers)
    clu.cmd.requires.setup_args(subparsers)

    return parser.parse_args(namespace=config)


def main() -> int:
    if sys.version_info < __about__.__minimum_python_info__:
        print(f"Must use at least python {__about__.__minimum_python__}", file=sys.stderr)
        sys.exit(1)

    #parser.parse_args(namespace=config)
    config = parse_cmdline()
    print(config)

    if config.verbosity is not logging.WARNING:
        logging.basicConfig(level=config.verbosity)

    log.info(f"Starting clu utility... {sys.argv=}")
    # if config.mode == "list-requires":
    #     do_list_requires()
    # elif config.mode == "check-requires":
    #     do_check_requires()
    # elif config.mode == "archive":
    #     do_archive()
    # elif config.mode == "report":
    #     do_report_facts()

    # command_map = {
    #     "archive": clu.cmd.archive.run,
    #     "report": clu.cmd.report.run,
    #     "requires": clu.cmd.requires.run,
    # }

    # Run the requested command
    # command_map[config.cmd]()
    print(config.cmd)
    config.func()
    return 0


if __name__ == "__main__":
    sys.exit(main())
