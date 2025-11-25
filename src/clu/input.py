"""Doc Incomplete."""

import logging
from pathlib import Path
import shlex
import shutil
import subprocess
from typing import Optional, Union

log = logging.getLogger(__name__)

FILE_SIZE_LIMIT = 1 * 1024 * 1024  # 1MB - arbitrary limit - seems reasonable


def text_file(fname: Union[str, Path], optional: bool = False) -> str:
    """Read a file and return its contents."""
    log.debug(f"text_file: {fname=}")

    fname = Path(fname)
    if not fname.is_file():
        if optional:
            log.info(f"Optional file not found: {fname}")
        else:
            log.error(f"File not found: {fname}")
        return ""

    file_size = fname.stat().st_size
    log.debug(f"File size is for {fname}: {file_size}")
    if file_size > FILE_SIZE_LIMIT:
        log.error(f"File size exceeds limit: {file_size} > {FILE_SIZE_LIMIT}")
        return ""

    # TODO: check permissions

    with open(fname, "r") as f:
        return f.read()


def text_program(cmdline) -> tuple[str, int]:
    log.debug(f"text_program: {cmdline}")

    # TODO: check for existence FIRST
    # if not check_program_exists(cmdline):
    #     log.error(f"Program not found: {cmdline}")
    #     return "", 1

    # TODO: check for permissions

    # TODO: check for output size

    try:
        result = subprocess.run(shlex.split(cmdline), capture_output=True, text=True)
        return result.stdout, result.returncode
    except Exception as e:
        log.error(f"Error running program {cmdline}: {e}")
        return "", 1


def transform_cmdline_to_filename(cmdline: str) -> tuple[str, str]:
    # cmdline is space separated, so we need to convert spaces to underscores
    cmdline = cmdline.replace(" ", "_")

    # udevadm info uses path like things that are not really paths - get rid of slashes
    cmdline = cmdline.replace("/", "%")

    log.debug(f"transformed {cmdline=}")
    return cmdline, cmdline + "_rc"


def check_program_exists(program: str) -> Optional[str]:
    log.debug(f"check_program_exists: {program}")

    return shutil.which(program.split()[0])  # Only check the actual command, not its arguments


def check_file_exists(fname: Union[str, Path]) -> Optional[str]:
    log.debug(f"check_file_exists: {fname}")

    fname = Path(fname)
    exists = fname.is_file()
    if exists:
        log.debug(f"File {fname} found")
        return str(fname)
    else:
        log.debug(f"File {fname} not found")
        return None
