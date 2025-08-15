"""Doc Incomplete."""

import logging
import shlex
import shutil
import subprocess
from typing import Optional
from pathlib import Path

log = logging.getLogger(__name__)


def text_file(fname) -> Optional[str]:
    """Read a file and return its contents.

    We [will] do checks for existance, size, perms and such thus the wrapper.
    """
    log.debug(f"text_file: {fname=}")

    if not Path(fname).is_file():
        log.info(f"File not found: {fname}")
        return None
    with open(fname, "r") as f:
        return f.read()


def transform_cmdline_to_filename(cmdline):
    log.debug(f"transform_cmdline_to_filename: {cmdline}")

    # cmdline is space separated, so we need to convert spaces to underscores
    cmdline = cmdline.replace(" ", "_")

    # udevadm info uses path like things that are not really paths - get rid of slashes
    cmdline = cmdline.replace("/", "%")

    log.debug(f"transformed {cmdline=}")
    return cmdline, cmdline + "_rc"


def text_program(cmdline) -> tuple[Optional[str], int]:
    log.debug(f"text_program: {cmdline}")

    try:
        result = subprocess.run(shlex.split(cmdline), capture_output=True, text=True)
        return result.stdout, result.returncode
    except Exception as e:
        log.debug(f"Error running program {cmdline}: {e}")
        return None, 1


def check_program_exists(program):
    log.debug(f"check_program_exists: {program}")

    return shutil.which(program.split()[0])  # Only check the actual command, not its arguments


def check_file_exists(fname):
    log.debug(f"check_file_exists: {fname}")

    exists = Path(fname).is_file()
    if exists:
        log.debug(f"File {fname} found")
        return fname
    else:
        log.debug(f"File {fname} not found")
        return None
