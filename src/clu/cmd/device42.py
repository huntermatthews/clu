import logging

from clu.config import get_config
from clu.device42_api import get_host_info
from clu.device42_utils import output_host_info

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
    log.info("Checking device status...")
    # Implement the check logic here
    return 0


def do_query():
    log.info(f"Querying device information... {cfg.hostname}")
    host_info = get_host_info(cfg.hostname)
    output_host_info(host_info)

    return 0


def do_update():
    log.info("Updating device information...")
    # Implement the update logic here
    return 0
