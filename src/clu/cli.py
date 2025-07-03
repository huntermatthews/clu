"""The main() function implementation for the program."""

import argparse
import sys
from pathlib import Path

# import clu
from clu import __about__, config
from clu.debug import debug, debug_var, trace, panic
from clu.report import do_report
from clu.archive import do_archive
from clu.requires import do_list_requires, do_check_requires

def main():
    if sys.version_info < (3, 8):
        print("ERROR: Must use at least python 3.8")
        sys.exit(1)

    parser = argparse.ArgumentParser(prog="clu", description="clu utility")
    parser.add_argument(
        "--version", action="version", version="%%(prog)s %s" % __about__.__version__
    )
    parser.add_argument(
        "--debug",
        action="count",
        default=0,
        help="Debug mode (use multiple times for more verbosity)",
    )
    parser.add_argument("--mock", help="Use mock data from a directory")
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

    if config.mock:
        # BUG: This is WEAK - we should use a more robust way to find the mock data directory
        config.mock = Path(__file__).parent.parent.parent / "mock_data" / config.mock
        if not config.mock.is_dir():
            panic(f"mock directory {config.mock} does not exist")
        debug_var("MOCK", config.mock)

    trace(f"Starting clu utility... {sys.argv=}")
    if config.mode == "list-requires":
        do_list_requires()
    elif config.mode == "check-requires":
        do_check_requires()
    elif config.mode == "archive":
        do_archive()
    elif config.mode == "report":
        do_report()

    return 0
