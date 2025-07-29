"""Doc Incomplete."""

import logging
import os
import shlex
import shutil
import subprocess


from clu import config

log = logging.getLogger(__name__)


def read_file(fname):
    log.debug(f"read_file: {fname=}")
    if config.mock:
        log.debug(f'config mocking {config.mock}')
        fname = get_file_mock_path(fname)
    data = raw_read_file(fname)
    return data


def raw_read_file(fname):
    if not os.path.isfile(fname):
        log.debug(f"File not found: {fname}")
        return None
    log.debug(f"Reading file: {fname}")
    with open(fname, "r") as f:
        return f.read()


def get_file_mock_path(fname):
    # FIX: stupid os.path.join screws up if any entry STARTS with a slash...
    x = os.path.join(config.mock, fname.strip('/'))
    log.debug(f"{x=}")
    return x


def transform_cmdline_to_filename(cmdline):
    log.debug(f"transform_cmdline_to_filename: {cmdline}")

    # cmdline is space separated, so we need to convert spaces to underscores
    cmdline = cmdline.replace(" ", "_")

    # udevadm info uses path like things that are not really paths - get rid of slashes
    cmdline = cmdline.replace("/", "%")

    log.debug(f"transformed {cmdline=}")
    return cmdline, cmdline + "_rc"


def get_program_mock_path(cmdline):
    cmd_name, rc_name = transform_cmdline_to_filename(cmdline)

    data_path = os.path.join(config.mock, "_programs", cmd_name)
    rc_path = os.path.join(config.mock, "_programs", rc_name)
    return (data_path, rc_path)


def read_program(cmdline):
    log.debug(f"read_program: {cmdline}")

    if config.mock:
        (dname, rc_name) = get_program_mock_path(cmdline)
        data = raw_read_file(dname)

        if os.path.isfile(rc_name):
            rc = int(raw_read_file(rc_name))
        elif data is None:
            rc = 127       # command not found (fish/bash/zsh/sh all consistent)
        else:
            rc = 0
        log.debug(f"{cmdline} {data=}")
        log.debug(f"{cmdline} {rc=}")
        return data, rc

    try:
        result = subprocess.run(shlex.split(cmdline), capture_output=True, text=True)
        return result.stdout, result.returncode
    except Exception as e:
        log.debug(f"Error running program {cmdline}: {e}")
        return None, 1


def check_program_exists(program):
    log.debug(f"check_program_exists: {program}")

    if config.mock:
        (dname, _) = get_program_mock_path(program)
        log.debug(f"mock path {dname=}")
        exists = os.path.isfile(dname)
        log.debug(f"{program} {exists=}")
        if exists:
            log.debug(f"Program {program} found in mock path: {dname}")
            return dname
        else:
            log.debug(f"Program {program} not found in mock path: {dname}")
            return None

    return shutil.which(program.split()[0])  # Only check the actual command, not its arguments


def check_file_exists(fname):
    log.debug(f"check_file_exists: {fname}")

    if config.mock:
        fname = get_file_mock_path(fname)

    exists = os.path.isfile(fname)
    if exists:
        log.debug(f"File {fname} found")
        return fname
    else:
        log.debug(f"File {fname} not found")
        return None
