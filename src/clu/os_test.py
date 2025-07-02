from clu.debug import trace, panic

from clu.os_linux import parse_cpuinfo_flags


def requires_os_test():
    trace("requires_os_test begin")
    panic("not implemented")


def parse_os_test():
    trace("parse_os_test begin")
    parse_cpuinfo_flags()
