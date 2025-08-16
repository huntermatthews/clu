import argparse
import json
import logging

from clu.opsys.factory import opsys_factory
from clu import Facts
from clu.opsys import OpSys


log = logging.getLogger(__name__)


def parse_args(subparsers):
    subp_report = subparsers.add_parser("report")
    subp_report.set_defaults(func=report_facts)

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


def set_report_defaults(opsys: OpSys, args: argparse.Namespace) -> None:
    """This is written defensively because if the user didn't say "report" explicitly on the command
    line, we want to make sure we still have a valid set of facts to report on.
    """
    if "output" not in args:
        args.output = "dots"

    if "all" in args and args.all:
        # BUG: this does not include the bmc facts..
        args.facts = "os sys phy run salt clu"
    elif "facts" not in args or not args.facts:
        args.facts = opsys.default_facts()


def parse_facts_by_specs(provides_map, parsed_facts: Facts, fact_specs) -> None:
    sources_to_parse = set()

    # Loop through the facts that were requested on the command line and get a set of parsers that
    # will obtain those facts (there's likely duplicates, so we use a set here)
    for fact_spec in fact_specs:
        for key in provides_map:
            if key.startswith(fact_spec):
                sources_to_parse.add(provides_map[key])

    # Call the parsers that we found in the previous loop
    for source in sources_to_parse:
        log.debug(f"Calling parser function {source}")
        source.parse(parsed_facts)


def filter_facts(requested_fact_specs, parsed_facts: Facts) -> Facts:
    """Filter the parsed facts based on the requested fact specifications."""

    # Loop through the facts that were requested on the command line (again) and make a new
    # facts list/dict that has JUST those matching facts
    output_facts = Facts()

    for fact_spec in requested_fact_specs:
        for key in parsed_facts:
            if key.startswith(fact_spec):
                output_facts[key] = parsed_facts[key]

    return output_facts


def do_output(output_facts: Facts, output_arg: str) -> None:
    if output_arg == "json":
        output_json(output_facts)
    elif output_arg == "shell":
        output_shell(output_facts)
    elif output_arg == "dots":
        output_dots(output_facts)


def report_facts(args) -> int:
    """Generate a report based on the current OS."""

    log.debug(f"Running command {args.cmd} with args={args}")

    opsys = opsys_factory()
    provides_map = opsys.provides()
    parsed_facts = Facts()

    set_report_defaults(opsys, args)

    parse_facts_by_specs(provides_map, parsed_facts, opsys.early_facts())
    parse_facts_by_specs(provides_map, parsed_facts, args.facts)

    output_facts = filter_facts(args.facts, parsed_facts)

    do_output(output_facts, args.output)

    return 0


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
