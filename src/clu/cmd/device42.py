import logging

from build.clu.opsys.factory import opsys_factory
from clu.config import get_config
from clu.device42.query import output_host_info, transform_host_info
from clu.device42.check import check_os_info
from clu.facts import Facts

log = logging.getLogger(__name__)
cfg = get_config()


def parse_args(subparsers):
    subp_device42 = subparsers.add_parser("device42", help="Interact with Device42.")
    subp_device42.set_defaults(func=do_run)

    device42_cmds = subp_device42.add_subparsers(
        dest="subcmd",
        description="Use 'clu device42 <command> --help' for more info.",
        # metavar="",  # setting this to empty string removes the ugly {foo,bar} from the help
        #        output
        title="subcommands",
    )

    d42_check = device42_cmds.add_parser("check", help="Check device status.")

    d42_query = device42_cmds.add_parser("query", help="Query device information.")
    d42_query.add_argument("hostname", help="Name of the host to query")

    d42_update = device42_cmds.add_parser("update", help="Update device information. Writes to db!")


def do_run():
    log.info(f"Running command {cfg.subcmd} with cfg={cfg}")
    if cfg.subcmd == "check":
        return do_check()
    elif cfg.subcmd == "query":
        return do_query()
    elif cfg.subcmd == "update":
        return do_update()
    else:
        raise ValueError(f"Unknown sub-command: {cfg.subcmd}")


def do_check():
    log.info(f"Begin checking {cfg.hostname} device information...")
    host_info = get_host_info(cfg.hostname)
    if host_info:
        opsys = opsys_factory()
        provides_map = opsys.provides()
        parsed_facts = Facts()

        check_os_info(host_info)
    else:
        log.error(f"Host {cfg.hostname} NOT found in Device42.")
        return 1
    log.info(f"End checking {cfg.hostname} device information...")
    return 0


def do_query():
    log.info(f"Querying device information... {cfg.hostname}")
    host_info = get_host_info(cfg.hostname)
    if host_info:
        print("Host Information:")
        transformed_info = transform_host_info(host_info)
        output_host_info(transformed_info)
        return 0
    else:
        print(f"No host information found for {cfg.hostname}.")
        return 1


def do_update():
    log.info("Updating device information...")
    # Implement the update logic here
    return 0
