"""Doc Incomplete."""

import os
import subprocess


from clu import facts
import clu
from clu.debug import trace, debug, debug_var, debug_var_list, trace_var_list, panic


def read_file(fname):
    trace("read_file begin")
    if clu.config.mock:
        fname = os.path.join(clu.config.mock, fname)
        debug_var("fname", fname)
    if not os.path.isfile(fname):
        debug(f"File not found: {fname}")
        return None
    debug(f"Reading file: {fname}")
    with open(fname, "r") as f:
        return f.read()


def read_program(p_name, *p_args):
    trace("read_program begin")
    arg_string = " ".join(p_args)
    debug_var("p_name", p_name)
    debug_var("p_args", p_args)
    debug_var("arg_string", arg_string)
    if clu.config.mock:
        fname = os.path.join(clu.config.mock, "_programs", p_name)
        if os.path.isfile(fname):
            with open(fname, "r") as f:
                return f.read()
        else:
            debug(f"Mock program file not found: {fname}")
            return None
    else:
        try:
            result = subprocess.run(
                [p_name] + list(p_args), capture_output=True, text=True, check=True
            )
            return result.stdout, result.returncode
        except Exception as e:
            debug(f"Error running program {p_name}: {e}")


def read_program2(*args):

    try:
        result = subprocess.run(args, capture_output=True, text=True, check=True)
        return result.stdout
    except Exception:
        return None
