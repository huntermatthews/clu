import os

from clu.debug import panic
from clu.opsys import OpSys
from clu.opsys.darwin import Darwin
# from clu.opsys.linux import Linux


class Unknown(OpSys):
    """Unknown operating system class."""

    # BUG: Implement parsing for unknown OS
    pass


def opsys_factory() -> OpSys:
    """Factory function to get the OS class based on the OS name."""

    os_name = os.uname().sysname
    if os_name == "Darwin":
        return Darwin()
    # elif os_name == "Linux":
    #     return Linux()
    else:
        return Unknown()
