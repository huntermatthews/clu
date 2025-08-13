import logging

from clu import Facts, Provides, Requires, panic

log = logging.getLogger(__name__)


def provides_macos_name(provides: Provides) -> None:
    provides["os.code_name"] = parse_macos_name


def requires_macos_name(requires: Requires) -> None:
    # its just logic code - there are no external requirements for this
    pass


def parse_macos_name(facts: Facts) -> None:
    if "os.version" not in facts:
        parse_sw_vers(facts)

    version = facts["os.version"]
    if not version:
        panic("parse_macos_name: os.version is not set or empty")

    major_ver = version.split(".")[0]
    log.debug(f"{major_ver=}")

    if major_ver == "26":
        code_name = "Tahoe"
    # BUG: check for 16-25 and error out
    elif major_ver == "15":
        code_name = "Sequoia"
    elif major_ver == "14":
        code_name = "Sonoma"
    elif major_ver == "13":
        code_name = "Ventura"
    elif major_ver == "12":
        code_name = "Monterey"
    elif major_ver == "11":
        code_name = "Big Sur"
    else:
        # Note that for older than 11, the logic of the code name changes
        # and thats WAY out of support for us
        code_name = f"Unknown-{major_ver}"
    facts["os.code_name"] = code_name
