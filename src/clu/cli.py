"""The main() function implementation for the program."""

import argparse
import sys
from pathlib import Path

import clu
from clu import __about__
from clu.debug import debug, debug_state, debug_var, trace, panic
from clu.report import do_report


def main():
    #    global clu.config
    print(f"{clu.config=}")

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

    clu.config = parser.parse_args()
    print(f"PDEBUG: clu.config = {clu.config}")
    print(f"PDEBUG: facts = {clu.config.facts}")

    if clu.config.debug == 1:
        debug_state("debug")
    elif clu.config.debug >= 2:
        debug_state("trace")
    else:
        debug_state("off")

    if clu.config.mock:
        # BUG: This is WEAK - we should use a more robust way to find the mock data directory
        clu.config.mock = (
            Path(__file__).parent.parent.parent / "mock_data" / clu.config.mock
        )
        if not clu.config.mock.is_dir():
            panic(f"mock directory {clu.config.mock_dir} does not exist")
        debug_var("MOCK", clu.config.mock)

    trace("Starting clu utility...")
    debug("Debugging enabled, debug state:", debug_state())
    debug_var("clu.config", clu.config)
    do_report()
    sys.exit(0)
