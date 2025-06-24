from clu.debug import trace, panic

def requires_os_unsupported():
    trace("requires_os_unsupported begin")
    panic("Unsupported OS")


def parse_os_unsupported():
    trace("parse_os_unsupported begin")
    panic("Unsupported OS")
