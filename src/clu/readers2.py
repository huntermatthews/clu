import os


def trace(msg):
    print(f"TRACE: {msg}")


def debug(msg):
    print(f"DEBUG: {msg}")


def debug_var(var_name, value):
    print(f"DEBUG: {var_name} = {value}")


MOCK = os.environ.get("MOCK", None)


def read_file(fname):
    trace("read_file begin")
    if MOCK:
        fname = os.path.join(MOCK, fname)
        debug_var("fname", fname)
    if not os.path.isfile(fname):
        debug(f"File not found: {fname}")
        return None
    debug(f"Reading file: {fname}")
    with open(fname, "r") as f:
        content = f.read()
    return content


import subprocess


def read_program(p_name, *p_args):
    trace("read_program begin")
    arg_string = " ".join(p_args)
    debug_var("p_name", p_name)
    debug_var("p_args", p_args)
    debug_var("arg_string", arg_string)
    if MOCK:
        fname = os.path.join(MOCK, "_programs", p_name)
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
            return result.stdout
        except Exception as e:
            debug(f"Error running program {p_name}: {e}")
