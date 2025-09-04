import json
import logging

from clu.opsys.factory import opsys_factory
from clu import Facts
from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()


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
    subp_report.add_argument(
        "--net",
        action="store_true",
        default=False,
        help="Enable network access - required for certain operations.",
    )
    subp_report.add_argument("facts", nargs="*", help="Facts to report on")


def set_report_defaults(default_facts: list[str], all_facts: list[str]) -> None:
    """If the user didn't say "report" explicitly on the command
    line, we want to make sure we still have a valid set of cfg to control our reporting.
    """

    # handle completely missing cfg
    if "output" not in cfg:
        cfg.output = "dots"
    if "all" not in cfg:
        cfg.all = False

    # Handle what the cfg mean
    if "facts" not in cfg or not cfg.facts:
        cfg.facts = default_facts
    elif "all" in cfg and cfg.all:
        cfg.facts = all_facts


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
        log.info(f"Calling parser function {source}")
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


def report_facts() -> int:
    """Generate a report based on the current OS."""

    log.info(f"Running command {cfg.cmd} with cfg={cfg}")

    opsys = opsys_factory()
    provides_map = opsys.provides()
    parsed_facts = Facts()

    set_report_defaults(opsys.default_facts(), list(provides_map.keys()))

    parse_facts_by_specs(provides_map, parsed_facts, opsys.early_facts())

    if cfg.all:
        cfg.facts = list(provides_map.keys())

    parse_facts_by_specs(provides_map, parsed_facts, cfg.facts)

    output_facts = filter_facts(cfg.facts, parsed_facts)

    do_output(output_facts, cfg.output)

    return 0


def output_dots(facts: Facts) -> None:
    for key in sorted(facts):
        value = facts[key]
        print(f"{key}: {value}")


def output_shell(facts: Facts) -> None:
    for key in sorted(facts):
        value = facts[key]
        key_var = key.upper().replace(".", "_")
        print(f'{key_var}="{value}"')


def output_json(facts: Facts) -> None:
    print(json.dumps(facts, indent=2))
