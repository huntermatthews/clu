"""Common code for all unit tests."""

from pathlib import Path

from clu.input import text_file, transform_cmdline_to_filename

mock_data_dir: Path = Path()


def get_file_mock_path(mock_dir, fname):
    """Get the mock file path for a given file name."""
    x = mock_dir / fname
    return x


def mock_read_file(mock_dir: Path, fname: Path) -> str:
    fname = get_file_mock_path(mock_dir, fname)
    data = text_file(fname)
    return data


def get_program_mock_path(mock_dir, cmdline):
    cmd_name, rc_name = transform_cmdline_to_filename(cmdline)

    data_path = mock_dir / "_programs" / cmd_name
    rc_path = mock_dir / "_programs" / rc_name
    return (data_path, rc_path)


def mock_read_program(mock_dir: Path, cmdline):
    """Read a program's output from the mock directory."""

    (dname, rc_name) = get_program_mock_path(mock_dir, cmdline)
    dname = Path(dname)
    rc_name = Path(rc_name)  # why pytest, why?

    data = text_file(Path(dname))

    # First is simple - no file, its command not found
    if not dname.is_file():
        rc = 127  # command not found (fish/bash/zsh/sh all consistent = 127)
    elif rc_name.is_file():
        # if we have a rc file, thats our RC by definition
        rc = int(text_file(Path(rc_name)))
    else:
        # otherwise, we can safely assume 0
        rc = 0

    return data, rc


def dict_subset(input_dict: dict, keys: list[str]) -> dict:
    """Return a new dictionary containing only the specified keys from the input dictionary."""
    # {k: input_dict.get(k, "") for k in ('l', 'm', 'n')}   # to synthesize missing keys

    return {key: input_dict[key] for key in keys if key in input_dict}


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
