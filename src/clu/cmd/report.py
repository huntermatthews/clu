import json
import logging

from clu.opsys.factory import opsys_factory
from clu import Facts


log = logging.getLogger(__name__)


def parse_args(subparsers):
    subp_report = subparsers.add_parser("report")
    subp_report.set_defaults(func=report_facts)

    subp_report.add_argument(
        "--test", action="store_true", help="Bypass uname checking and do whatever"
    )
    subp_report.add_argument(
        "--output",
        choices=["dots", "shell", "json"],
        default="dots",
        help="Output format: 'dots', 'shell', or 'json'",
    )
    subp_report.add_argument(
        "--all",
        "-A",
        action="store_true",
        help="Output all facts",
    )
    subp_report.add_argument("facts", nargs="*", help="Facts to report on")


def report_facts(args) -> None:
    """Generate a report based on the current OS."""

    log.debug(f"Running command {args.cmd} with args={args}")

    opsys = opsys_factory()

    provides_map = opsys.provides()
    parsers_to_call = set()

    if args.all:
        args.facts = "os sys phy run salt clu"
    elif not args.facts:
        args.facts = opsys.default_facts()

    # Loop through the facts that were requested on the command line and get a set of parsers that
    # will obtain those facts (there's likely duplicates, so we use a set here)
    for fact_spec in args.facts:
        for key in provides_map:
            if key.startswith(fact_spec):
                parsers_to_call.add(provides_map[key])

    # Call the parsers that we found in the previous loop
    parsed_facts = Facts()
    for parser_fn in parsers_to_call:
        parser_fn(parsed_facts)

    # Loop through the facts that were requested on the command line (again) and make a new
    # facts list/dict that has JUST those matching facts
    output_facts = Facts()
    for fact_spec in args.facts:
        for key in parsed_facts:
            if key.startswith(fact_spec):
                output_facts[key] = parsed_facts[key]

    if args.output == "json":
        output_json(output_facts)
    elif args.output == "shell":
        output_shell(output_facts)
    else:
        output_dots(output_facts)


def output_dots(facts: Facts) -> None:
    for key in facts:
        value = facts[key]
        print(f"{key}: {value}")


def output_shell(facts: Facts) -> None:
    for key in facts:
        value = facts[key]
        key_var = key.upper().replace(".", "_")
        print(f'{key_var}="{value}"')


def output_json(facts: Facts) -> None:
    print(json.dumps(facts, indent=2))
