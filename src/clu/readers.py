"""Doc Incomplete."""

import os
import shlex
import subprocess


from clu import config
from clu.debug import trace, debug, debug_var, debug_var_list, trace_var_list, panic


def read_file(fname):
    trace("read_file begin")
    debug(f"read_file: {fname=}")
    if config.mock:
        fname = text_file_mock_path(fname)
    data = raw_read_file(fname)
    return data

def raw_read_file(fname):
    trace("raw_read_file begin")
    if not os.path.isfile(fname):
        debug(f"File not found: {fname}")
        return None
    debug(f"Reading file: {fname}")
    with open(fname, "r") as f:
        return f.read()

def text_file_mock_path(fname):
    return os.path.join(config.mock, fname)

def program_mock_path(cmdline):
    trace("program_mock_path begin")

    # cmdline is space separated, so we need to convert spaces to underscores
    cmdline = cmdline.replace(" ", "_")
    data_file = os.path.join(config.mock, "_programs", cmdline)
    rc_file = data_file + ".rc"
    return (data_file, rc_file)


def read_program(cmdline):
    trace("read_program begin")
    debug(f"read_program: {cmdline}")

    if config.mock:
        (dname, rc_name) = program_mock_path(cmdline)
        data = raw_read_file(dname)

        rc = 0   # Default return code
        if os.path.isfile(rc_name):
            rc = int(raw_read_file(rc_name))
        return data, rc

    try:
        result = subprocess.run(shlex.split(cmdline), capture_output=True, text=True)
        return result.stdout, result.returncode
    except Exception as e:
        debug(f"Error running program {cmdline}: {e}")
        return None, 1


def read_program2(*args):

    try:
        result = subprocess.run(args, capture_output=True, text=True, check=True)
        return result.stdout
    except Exception:
        return None


def _subprocess_run(cmdline):
    # TODO: Should this simply be merged into read_simple_program() ?
    result = subprocess.run(
        shlex.split(cmdline),
        check=False,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        encoding="UTF-8",
    )
    return result
