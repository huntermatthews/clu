import logging

from clu.config import get_config
from clu.device42.query import subcmd_query
from clu.device42.check import subcmd_check


log = logging.getLogger(__name__)
cfg = get_config()


def parse_args(subparsers):
    subp_device42 = subparsers.add_parser("device42", help="Interact with Device42.")

    device42_cmds = subp_device42.add_subparsers(
        dest="subcmd",
        description="Use 'clu device42 <sub-command> --help' for more info.",
        metavar="",  # setting this to empty string removes the ugly {foo,bar} from the help output
        title="subcommands",
    )

    d42_check = device42_cmds.add_parser("check", help="Check device status.")
    d42_check.set_defaults(func=subcmd_check)

    d42_query = device42_cmds.add_parser("query", help="Query device information.")
    # TODO: add default of this host
    d42_query.add_argument("hostname", help="Name of the host to query for")
    d42_query.set_defaults(func=subcmd_query)

    d42_update = device42_cmds.add_parser("update", help="Update device information. Writes to db!")
    d42_update.set_defaults(func=subcmd_update)


def subcmd_update():
    log.info("Not Implemented Yet")
    # # Implement the update logic here
    return 1
