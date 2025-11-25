import json
import logging

from clu.opsys.factory import opsys_factory
from clu.facts import Facts, Tier
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
        "--net",
        action="store_true",
        default=False,
        help="Enable network access - required for certain operations.",
    )
    subp_report.add_argument(
        "facts", nargs="*", help="Specfic facts to report on, regardless of tier."
    )

    ex_group = subp_report.add_mutually_exclusive_group()
    ex_group.add_argument(
        "-1",
        action="store_const",
        const=1,
        dest="tier",
        default=1,
        help="Report all tier 1 facts (the default, basic info only).",
    )
    ex_group.add_argument(
        "-2",
        action="store_const",
        const=2,
        dest="tier",
        help="Report all tier 1 And 2 facts (most of the details).",
    )
    ex_group.add_argument(
        "-3",
        action="store_const",
        const=3,
        dest="tier",
        help="Report all tier 1, 2, and 3 facts (everything, useless details included).",
    )


def set_report_defaults(all_facts: list[str]) -> None:
    """If the user didn't say "report" explicitly on the command
    line, we want to make sure we still have a valid set of configs to control our reporting.
    """

    # handle completely missing config items
    if "output" not in cfg:
        cfg.output = "dots"
    if "verbose_level" not in cfg:
        cfg.verbose_level = 0
    if "net" not in cfg:
        cfg.net = False
    if "tier" not in cfg:
        cfg.tier = 1

    if "facts" not in cfg or cfg.facts == []:
        cfg.facts = all_facts


def parse_facts_by_specs(provides_map, parsed_facts: Facts, fact_specs) -> None:
    sources_to_parse = set()

    # Loop through the facts that were requested on the command line and get a set of parsers that
    # will obtain those facts (there's likely duplicates for the sources, so we use a set here)
    if fact_specs:
        for fact_spec in fact_specs:
            for key in provides_map:
                if key.startswith(fact_spec):
                    sources_to_parse.add(provides_map[key])
    else:
        # If the user didn't request any specific facts, we need to parse all the sources
        sources_to_parse = set(provides_map.values())

    # Call the parsers that we found in the previous loop
    for source in sources_to_parse:
        log.info(f"Calling parser function {source}")
        source.parse(parsed_facts)


def filter_facts(parsed_facts: Facts, requested_fact_specs, tier) -> Facts:
    """Filter the parsed facts based on the requested fact specifications."""

    # Loop through the facts that were requested on the command line (again) and make a new
    # facts list/dict that has JUST those matching facts, taking into account the tier level too.
    output_facts = Facts()
    tier_facts = parsed_facts.get_tier(Tier.get_by_int(tier))

    for fact_spec in requested_fact_specs:
        for key in parsed_facts:
            if key.startswith(fact_spec) and key in tier_facts:
                # print(f"Adding fact {key} to output")
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

    # User may not have explicity said "report" as the command name on command line - fill in the
    # gaps of our config if so.
    set_report_defaults(provides_map.keys())
    log.info(f"Updated config cfg={cfg}")

    # Our crude inter-fact dependency code is just to hard-code which facts _might_ be depended on
    # by other sources.  We parse the sources that give us those facts here.
    parse_facts_by_specs(provides_map, parsed_facts, opsys.early_facts())

    # Now that we have all the early facts, using the list of fact names the user requested or a
    # tier, parse the sources for that list - this is main parsing step of the report command.
    parse_facts_by_specs(provides_map, parsed_facts, cfg.facts)

    # the sources always parse all the facts it can out of a file/program/whatever.
    # (its generally just as fast and FAR simpler for the Sources NOT to care)
    # But the user might have requested less than that - filter out the extra stuff here.
    output_facts = filter_facts(parsed_facts, cfg.facts, cfg.tier)

    # Finally, out the filtered facts in the format the user requested (might be defaulted)
    do_output(output_facts, cfg.output)

    # the sub-cmd return code will be the programs return code
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
    print(json.dumps(facts.to_dict(), indent=2))
