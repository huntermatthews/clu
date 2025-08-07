from clu.os_map import get_os_functions
from clu.readers import check_file_exists, check_program_exists
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

    (requires_fn, _, _, _) = get_os_functions()
    requires = requires_fn()

    requires.sort()
    return requires
