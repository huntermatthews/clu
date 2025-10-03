import logging

from clu.facts import Facts
from clu.config import get_config
from clu.opsys.factory import opsys_factory

log = logging.getLogger(__name__)
cfg = get_config()


# BUG: duplicated from cmds/report.py
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


def get_host_info(hostname: str) -> dict | None:
    opsys = opsys_factory()
    provides_map = opsys.provides()
    parsed_facts = Facts()


def check_os_info(host_info):
    log.info(f"Checking OS information for {host_info['name']}...")

    # We need the OS distro version from Sources

    log.info(f"OS information check complete for {host_info['name']}.")
