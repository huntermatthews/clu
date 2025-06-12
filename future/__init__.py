from pathlib import Path
import subprocess
import shlex
#import str

def debug(*args, **kwargs):
    if g_config.debug:
        return builtins.print(*args, **kwargs)

def verbose(*args, **kwargs):
    if g_config.verbose:
        return builtins.print(*args, **kwargs)

# BUG: should not print anything...
# a simple file is ONE line that contains ONE value
def read_simple_file(filename):
    data = None
    pname = Path(filename)
    try:
        with pname.open() as f:
            data = f.read().strip()
    except IOError as ex:
        print(ex.errno, "XX", ex.strerror, "XX", ex.filename)
    return data


def read_simple_program(cmdline):
    # This cheesy wrapper should only be used on "small" things - you've been warned.
    # BUG: obviously we should verify the program is there (exist)
    # BUG: we can excute it (perms)
    try:
        result = _subprocess_run(cmdline)
    except FileNotFoundError as ex:
        # BUG: log this
        return None
    if result.returncode != 0:
        return None
    return result.stdout

def _subprocess_run(cmdline):
    # TODO: Should this simply be merged into read_simple_program() ?
    result = subprocess.run(shlex.split(cmdline),
                            check=False,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE,
                            encoding='UTF-8')
    return result


def simple_parser(text, seperator="."):
    # Just what it says on the tin - do a VERY rough parsing of multiple lines
    # thing on the left of the Seperator, key
    # thing on the right of the seperator, value
    # multiple seperators? ignored -only first  one counts
    # no key? NO PROBLEM - ignore the whole line!
    # redundant keys? NO PROBLEM - last value wins!
    # its .... simple
    data = {}

    for line in text.splitlines():
        if not line.strip():
            continue
        if seperator not in line:
            continue
        key, value = line.split(seperator, 1)
        data[key.strip()] = value.strip()

    return data

    ## END OF LINE ##


