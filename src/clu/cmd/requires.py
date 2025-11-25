import logging

from clu.config import get_config
from clu.input import check_file_exists, check_program_exists
from clu.opsys.factory import opsys_factory

log = logging.getLogger(__name__)
cfg = get_config()


def parse_args(subparsers):
    subp_requires = subparsers.add_parser(
        "requires", help="Check or list requirements for the system."
    )
    subp_requires.set_defaults(func=run)

    subp_requires.add_argument("subcmd", choices=["list", "check"], help="Sub-Command to run")


def run():
    log.info(f"Running command {cfg.cmd} with cfg={cfg}")
    if cfg.subcmd == "list":
        return list_requires()
    elif cfg.subcmd == "check":
        return check_requires()
    else:
        raise ValueError(f"Unknown sub-command: {cfg.subcmd}")


def list_requires() -> int:
    """List all the requirements for the current OS."""

    requires = opsys_factory().requires()

    print("Listing Requirements: ")
    print("----------------------")

    print("Files:")
    for file in requires.files:
        print(f"  - {file}")

    print("Programs:")
    for program in requires.programs:
        print(f"  - {program}")

    # print("APIs:")
    # for api in requires.apis:
    #     print(f"  - {api}")

    return 0


def check_requires() -> int:
    """Check if the current OS meets the requirements."""

    requires = opsys_factory().requires()
    print("Checking requirements:")
    print("----------------------")

    print("Files:")
    for file in requires.files:
        if check_file_exists(file):
            print(f"  - [ok] {file}")
        else:
            print(f"  - [MISSING] {file}")

    print("Programs:")
    for program in requires.programs:
        if check_program_exists(program):
            print(f"  - [ok] {program}")
        else:
            print(f"  - [MISSING] {program}")

    # print("APIs:")
    # for api in requires["apis"]:
    #     # what would go here?
    #     pass

    return 0
