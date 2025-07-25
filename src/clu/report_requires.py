from clu import config
from clu.facts import Facts
from clu.os_darwin import requires_os_darwin
from clu.os_linux import requires_os_linux
from clu.os_unsupported import requires_os_unsupported
from clu.os_generic import parse_uname
from clu.readers import check_file_exists, check_program_exists
from clu.os_test import requires_os_test
from clu.requires import Requires

def do_list_requires() -> None:
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


def do_check_requires() -> None:
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

    if config.test:
    # If we're in test mode, we don't need to do any checks.
    # We just requires the test OS.
        requires = Requires()
        requires_os_test(requires)
    else:
        # we always require uname, so we parse it first
        facts = Facts()
        parse_uname(facts)

        if facts["os.kernel.name"] == "Darwin":
            requires_os_darwin(requires)
        elif facts["os.kernel.name"] == "Linux":
            requires_os_linux(requires)
        else:
            requires_os_unsupported(requires)

    requires.sort()
    return requires
