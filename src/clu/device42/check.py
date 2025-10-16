import logging

from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()


def subcmd_check():
    log.info("Begin checking this hosts information...")
    # host_info = get_host_info(cfg.hostname)
    # if host_info:
    #     opsys = opsys_factory()
    #     provides_map = opsys.provides()
    #     parsed_facts = Facts()

    #     check_os_info(host_info)
    # else:
    #     log.error(f"Host {cfg.hostname} NOT found in Device42.")
    #     return 1
    # log.info(f"End checking {cfg.hostname} device information...")
    return 55


def check_os_info(host_info):
    log.info(f"Checking OS information for {host_info['name']}...")

    # We need the OS distro version from Sources

    log.info(f"OS information check complete for {host_info['name']}.")
