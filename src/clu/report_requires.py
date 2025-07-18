from clu.facts import get_fact
from clu.os_darwin import requires_os_darwin
from clu.os_linux import requires_os_linux
from clu.os_unsupported import requires_os_unsupported
from clu.os_generic import parse_uname   # we always need this, to know WHICH requires lists to handle.
from clu.requires import get_all_requires
from clu.readers import check_file_exists, check_program_exists


def do_list_requires():
    """List all the requirements for the current OS."""

    get_os_requirements()

    print(f"Requirements: ({get_fact('os.kernel.name')})")

    print("Files:")
    for file in get_all_requires()["files"]:
        print(f"  - {file}")

    print("Programs:")
    for program in get_all_requires()["programs"]:
        print(f"  - {program}")

    print("APIs:")
    for api in get_all_requires()["apis"]:
         print(f"  - {api}")


def do_check_requires():
    """Check if the current OS meets the requirements."""

    get_os_requirements()

    print(f"Checking requirements: ({get_fact('os.kernel.name')})")

    print("Files:")
    for file in get_all_requires()["files"]:
        # BUG: we run as root, so we theoretically can access all files
        # but thats NOT always the case - see selinux and weird macos permissions
        if check_file_exists(file):
            print(f"  - [ok] {file}")
        else:
            print(f"  - [MISSING] {file}")

    print("Programs:")
    for program in get_all_requires()["programs"]:
        if check_program_exists(program):
            print(f"  - [ok] {program}")
        else:
            print(f"  - [MISSING] {program}")

    # print("APIs:")
    # for api in requires["apis"]:
    #     # what would go here?
    #     pass


def get_os_requirements():
    """Get the requirements for the current OS."""
    # we always require uname, so we parse it first
    parse_uname()

    if get_fact("os.kernel.name") == "Darwin":
        requires_os_darwin()
    elif get_fact("os.kernel.name") == "Linux":
        requires_os_linux()
    else:
        requires_os_unsupported()

    get_all_requires()["programs"].sort()
    get_all_requires()["files"].sort()
    get_all_requires()["apis"].sort()
