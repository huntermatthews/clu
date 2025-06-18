# -*- coding: utf-8 -*-

import argparse
import os
import sys
import platform
from pathlib import Path

_PROGRAM = os.path.basename(sys.argv[0])
_VERSION = "1.2"
_HELP = f"""
Usage: {_PROGRAM} [OPTIONS] [FILE]
Options:
  -h, --help      Show this help message and exit
  -v, --version   Show version information and exit
    other docs later
"""


def panic(msg):
    print(f"Error: {msg}", file=sys.stderr)
    sys.exit(1)


def debug_var(var_name, value):
    print(f"DEBUG: {var_name} = {value}")


def trace(msg):
    print(f"TRACE: {msg}")


def collector():
    print("Collector called (stub).")


def check_requires():
    print("Check requires called (stub).")


def list_requires():
    print("List requires called (stub).")


def os_darwin_parse():
    print("Darwin parse called (stub).")


def os_test_parse():
    print("Test parse called (stub).")


def os_unsupported_parse():
    print("Unsupported OS parse called (stub).")


def output_dots():
    print("Output: dots (stub).")


def output_shell():
    print("Output: shell (stub).")


def main():
    parser = argparse.ArgumentParser(
        prog=_PROGRAM, description="Gru utility", add_help=False, usage=_HELP
    )
    parser.add_argument("-h", "--help", action="store_true")
    parser.add_argument("-v", "--version", action="store_true")
    parser.add_argument("--debug", nargs="*")
    parser.add_argument("--mock")
    parser.add_argument("--output")
    parser.add_argument("-R", "--check-requires", action="store_true")
    parser.add_argument("-L", "--list-requires", action="store_true")
    parser.add_argument("--collector", action="store_true")
    parser.add_argument("file", nargs="?")

    args = parser.parse_args()

    if args.help:
        print(_HELP)
        sys.exit(0)

    if args.version:
        print(f"{_PROGRAM} v{_VERSION}")
        sys.exit(0)

    if args.debug is not None:
        if len(args.debug) > 1:
            _debug = "trace"
        else:
            _debug = "debug"
        debug_var("_debug", _debug)

    if args.mock:
        mock_dir = Path(__file__).parent / "mock_data" / args.mock
        if not mock_dir.is_dir():
            panic(f"{args.mock} directory {mock_dir} does not exist")
        debug_var("MOCK", mock_dir)

    if args.collector:
        trace("Calling Collector...")
        collector()
        sys.exit(0)

    kernel_name = platform.system()
    debug_var("kernel_name", kernel_name)

    if args.check_requires:
        trace("Calling check_requires...")
        check_requires()
        sys.exit(0)

    if args.list_requires:
        trace("Calling list_requires...")
        list_requires()
        sys.exit(0)

    if kernel_name == "Darwin":
        os_darwin_parse()
    elif kernel_name == "Linux":
        # os_linux_parse()
        os_test_parse()
    else:
        os_unsupported_parse()

    output = args.output or "dots"
    if output == "dots":
        output_dots()
    elif output == "shell":
        output_shell()
    else:
        output_dots()


if __name__ == "__main__":
    main()
