import logging

from clu.os_map import get_os_functions
from clu.input import check_file_exists, check_program_exists
from clu import Requires


log = logging.getLogger(__name__)


def parse_args(subparsers):
    subp_requires = subparsers.add_parser("requires")
    subp_requires.set_defaults(func=run)

    subp_requires.add_argument("subcmd", choices=["list", "check"], help="Sub-Command to run")


def run(args):
    log.debug(f"Running command {args.cmd} with args={args}")
    if args.subcmd == "list":
        list_requires()
    elif args.subcmd == "check":
        check_requires()
    else:
        raise ValueError(f"Unknown sub-command: {args.subcmd}")

    return 0

def list_requires() -> None:
    """List all the requirements for the current OS."""

    requires = get_os_requirements()

    print("Requirements: ")

    print("Files:")
    for file in requires.files:
        print(f"  - {file}")

    print("Programs:")
    for program in requires.programs:
        print(f"  - {program}")

    print("APIs:")
    for api in requires.apis:
        print(f"  - {api}")


def check_requires() -> None:
    """Check if the current OS meets the requirements."""

    requires = get_os_requirements()
    print("Checking requirements:")

    print("Files:")
    for file in requires.files:
        # BUG: we run as root, so we theoretically can access all files
        # but thats NOT always the case - see selinux and weird macos permissions
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


def get_os_requirements() -> Requires:
    """Get the requirements for the current OS."""

    (requires_fn, _, _, _) = get_os_functions()
    requires = requires_fn()

    requires.sort()
    return requires
