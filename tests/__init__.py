"""Common code for all unit tests."""

import os
from pathlib import Path

from clu.readers import read_file, transform_cmdline_to_filename


def get_file_mock_path(mock_dir, fname):
    """Get the mock file path for a given file name."""
    x = mock_dir / fname
    return x


def mock_read_file(mock_dir: Path, fname: str) -> str:
    fname = get_file_mock_path(mock_dir, fname)
    data = read_file(fname)
    return data


def get_program_mock_path(mock_dir, cmdline):
    cmd_name, rc_name = transform_cmdline_to_filename(cmdline)

    data_path = os.path.join(mock_dir, "_programs", cmd_name)
    rc_path = os.path.join(mock_dir, "_programs", rc_name)
    return (data_path, rc_path)


def mock_read_program(mock_dir: Path, cmdline):
    """Read a program's output from the mock directory."""

    (dname, rc_name) = get_program_mock_path(mock_dir, cmdline)
    data = read_file(dname)

    if os.path.isfile(rc_name):
        rc = int(read_file(rc_name))
    elif data is None:
        rc = 127  # command not found (fish/bash/zsh/sh all consistent = 127)
    else:
        rc = 0

    return data, rc


# For later tests

# def check_program_exists(program):
#     debug(f"check_program_exists: {program}")

#     if config.mock:
#         (dname, _) = get_program_mock_path(program)
#         debug_var("mock path", dname)
#         exists = os.path.isfile(dname)
#         debug_var(f"{program} exists", exists)
#         if exists:
#             debug(f"Program {program} found in mock path: {dname}")
#             return dname
#         else:
#             debug(f"Program {program} not found in mock path: {dname}")
#             return None

#     return shutil.which(program.split()[0])  # Only check the actual command, not its arguments

# def check_file_exists(fname):
#     debug(f"check_file_exists: {fname}")

#     if config.mock:
#         fname = get_file_mock_path(fname)

#     exists = os.path.isfile(fname)
#     if exists:
#         debug(f"File {fname} found")
#         return fname
#     else:
#         debug(f"File {fname} not found")
#         return None
