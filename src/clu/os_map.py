from clu.debug import panic
from clu import Facts
from clu.darwin import Darwin


def OsFactory() -> Darwin | None:
    """Factory function to get the OS class based on the OS name."""
    if os_name == "Darwin":
        return Darwin()
    elif os_name == "Linux":
        return Linux()
    return None
