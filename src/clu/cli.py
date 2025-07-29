"""The main() function implementation for the program."""

import argparse
import logging
import sys
from pathlib import Path

# import clu
from clu import __about__, config
from clu.report_facts import do_report_facts
from clu.archive import do_archive
from clu.report_requires import do_list_requires, do_check_requires

log = logging.getLogger(__name__)

def main()  -> int:
    if sys.version_info < (3, 8):
        log.error("Must use at least python 3.8")
        sys.exit(1)

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
    parser.add_argument("--test", action="store_true", help="Bypass uname checking and do whatever")
    parser.add_argument(
        "--output",
        choices=["dots", "shell"],
        default="dots",
        help="Output format: 'dots' or 'shell'",
    )
    group = parser.add_mutually_exclusive_group()
    group.add_argument(
        "--report",
        action="store_const",
        const="report",
        default="report",
        dest="mode",
    )
    group.add_argument(
        "--archive",
        action="store_const",
        const="archive",
        dest="mode",
    )
    group.add_argument(
        "--list-requires",
        action="store_const",
        const="list-requires",
        dest="mode",
    )
    group.add_argument(
        "--check-requires",
        action="store_const",
        const="check-requires",
        dest="mode",
    )

    parser.add_argument("facts", nargs="*", help="Facts to collect")

    parser.parse_args(namespace=config)

    if config.verbosity is not logging.WARNING:
        logging.basicConfig(level=config.verbosity)

    log.info(f"Starting clu utility... {sys.argv=}")
    if config.mode == "list-requires":
        do_list_requires()
    elif config.mode == "check-requires":
        do_check_requires()
    elif config.mode == "archive":
        do_archive()
    elif config.mode == "report":
        do_report_facts()

    return 0
