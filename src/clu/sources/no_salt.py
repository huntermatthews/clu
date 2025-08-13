import logging

from clu import Facts, Provides, Requires
from clu.input import text_file

log = logging.getLogger(__name__)


def provides_no_salt(provides: Provides) -> None:
    provides["os.no_salt.exists"] = parse_no_salt
    provides["os.no_salt.reason"] = parse_no_salt


def requires_no_salt(requires: Requires) -> None:
    requires.files.append("/no_salt")


def parse_no_salt(facts: Facts) -> None:
    data = text_file("/no_salt")
    if data is None:
        facts["salt.no_salt.exists"] = "False"
        return
    else:
        log.debug(f"{data=}")
        facts["salt.no_salt.exists"] = "True"
    if not data.strip():
        data = "UNKNOWN"
    facts["salt.no_salt.reason"] = data.strip()
