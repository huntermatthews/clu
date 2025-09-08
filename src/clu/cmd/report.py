import json
import logging
from typing import List

from clu.opsys.factory import opsys_factory
from clu import facts
from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()


def parse_args(subparsers):
    subp_report = subparsers.add_parser("report", help="Report facts about this system")
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


def set_report_defaults(default_facts: list[str]) -> None:
    """If the user didn't say "report" explicitly on the command
    line, we want to make sure we still have a valid set of configs to control our reporting.
    """

    # handle completely missing config
    # print('handling completely missing config', cfg)
    if "output" not in cfg:
        cfg.output = "dots"
    if "all" not in cfg:
        cfg.all = False

    if "facts" not in cfg or cfg.facts == []:
        # print('setting default facts')
        cfg.facts = default_facts


def parse_sources_by_fact_names(provides_map, requested_fact_names) -> None:
    sources_to_parse = set()

    # Loop through the fact names that were requested on the command line and get a set of
    # parsers that will obtain those facts (there's likely duplicates, so we use a set here)
    #    breakpoint()
    for fact_name in requested_fact_names:
        for key in provides_map:
            if key.startswith(fact_name):
                sources_to_parse.add(provides_map[key])

    # Call the parsers that we found in the previous loop
    for source in sources_to_parse:
        log.info(f"Calling parser function {source}")
        source.parse()


def filter_facts(requested_fact_names) -> List[str]:
    """Filter the parsed facts based on the requested fact specifications."""

    # Loop through the facts that were requested on the command line (again) and make a new
    # facts list/dict that has JUST those matching facts
    output_fact_names = []

    for fact_name in requested_fact_names:
        for name in facts:
            if name.startswith(fact_name):
                output_fact_names.append(name)

    return output_fact_names


def do_output(output_fact_names, output_arg: str) -> None:
    if output_arg == "json":
        output_json(output_fact_names)
    elif output_arg == "shell":
        output_shell(output_fact_names)
    elif output_arg == "dots":
        output_dots(output_fact_names)


def report_facts() -> int:
    """Generate a report based on the current OS."""

    log.info(f"Running command {cfg.cmd} with cfg={cfg}")

    # Get the correct provides map and default facts for the current OS
    opsys = opsys_factory()
    provides_map = opsys.provides()
    set_report_defaults(opsys.default_facts())

    # Our crude inter-fact dependency code is just to hard-code which facts _might_ be depended on
    # by other sources.  We parse the sources that give us those facts here.
    parse_sources_by_fact_names(provides_map, opsys.early_facts())

    if cfg.all:
        cfg.facts = list(provides_map.keys())

    # Now that we have all the early facts, using the list of fact names the user requested (might
    # be defaulted) parse the sources for that list - this is main parsing loop of the sub-command.
    parse_sources_by_fact_names(provides_map, cfg.facts)

    # the sources always parse all the facts it can out of a file/program/whatever.
    # (its generally just as fast and FAR simpler for the Sources NOT to care)
    # But the user might have requested less than that - filter out the extra stuff here.
    output_facts = filter_facts(cfg.facts)

    # Finally, out the filtered facts in the format the user requested (might be defaulted)
    do_output(output_facts, cfg.output)

    # the sub-cmd return code will be the programs return code
    return 0


def output_dots(output_fact_names: list[str]) -> None:
    for name in sorted(output_fact_names):
        value = facts[name]
        print(f"{name}: {value}")


def output_shell(output_fact_names: list[str]) -> None:
    for name in sorted(output_fact_names):
        value = facts[name]
        key_var = name.upper().replace(".", "_")
        print(f'{key_var}="{value}"')


def output_json(output_fact_names: list[str]) -> None:
    partial_copy = {name: facts[name] for name in output_fact_names if name in facts}
    print(json.dumps(partial_copy, indent=2))
