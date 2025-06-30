"""Doc Incomplete."""

import os
import shlex
import subprocess


from clu import config
from clu.debug import trace, debug, debug_var


def read_file(fname):
    trace("read_file begin")
    debug(f"read_file: {fname=}")
    if config.mock:
        debug(f'config mocking {config.mock}')
        fname = get_file_mock_path(fname)
    data = raw_read_file(fname)
    return data


def raw_read_file(fname):
    trace(f"raw_read_file begin {fname=}")
    if not os.path.isfile(fname):
        debug(f"File not found: {fname}")
        return None
    debug(f"Reading file: {fname}")
    with open(fname, "r") as f:
        return f.read()


def get_file_mock_path(fname):
    x = os.path.join(config.mock, fname.strip('/'))
    debug_var("x", x)
    return x


def get_program_mock_path(cmdline):
    trace("program_mock_path begin")

    # cmdline is space separated, so we need to convert spaces to underscores
    cmdline = cmdline.replace(" ", "_")

    # udevadm info uses path like things that are not really paths - get rid of slashes
    cmdline = cmdline.replace("/", "%")

    data_file = os.path.join(config.mock, "_programs", cmdline)
    rc_file = data_file + ".rc"
    return (data_file, rc_file)


def read_program(cmdline):
    trace("read_program begin")
    debug(f"read_program: {cmdline}")

    if config.mock:
        (dname, rc_name) = get_program_mock_path(cmdline)
        data = raw_read_file(dname)

        rc = 0  # Default return code
        if os.path.isfile(rc_name):
            rc = int(raw_read_file(rc_name))
        return data, rc

    try:
        result = subprocess.run(shlex.split(cmdline), capture_output=True, text=True)
        return result.stdout, result.returncode
    except Exception as e:
        debug(f"Error running program {cmdline}: {e}")
        return None, 1
