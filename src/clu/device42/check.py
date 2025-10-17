import logging
import os

from clu.config import get_config
from clu.device42.api import get_host_by_name
from clu.cmd.report import parse_facts_by_specs
from clu.facts import Facts
from clu.opsys.factory import opsys_factory

log = logging.getLogger(__name__)
cfg = get_config()


def subcmd_check():
    # we need our hostname
    # Can't be from CLI - we can't check any host but ourselves
    nodename = os.uname().nodename.split(".")[0]
    log.info(f"Querying host information... {nodename}")

    host_info = get_host_by_name(nodename)
    if not host_info:
        print(f"No host information found for {nodename}.")
        return 1
    else:
        print("Host Information:")
        opsys = opsys_factory()
        provides_map = opsys.provides()
        parsed_facts = Facts()

        parse_facts_by_specs(provides_map, parsed_facts, fact_specs)


def check_host_info(host_info):
    log.info(f"Checking host information for {host_info['name']}...")

    # We need the OS distro version from Sources

    log.info(f"Host information check complete for {host_info['name']}.")


def check_os_info(host_info):
    log.info(f"Checking OS information for {host_info['name']}...")

    # We need the OS distro version from Sources

    log.info(f"OS information check complete for {host_info['name']}.")
