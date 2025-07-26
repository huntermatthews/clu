""" Code common to all tests.
"""
from functools import partial
from pathlib import Path

import pytest


from clu.debug import debug
from clu.readers import raw_read_file, get_program_mock_path, read_program


@pytest.hookimpl
def pytest_configure(config: pytest.Config) -> None:
    pytest.mock_dir = config.rootdir / "mock_data"


def get_file_mock_path(mock_dir, fname):
    """Get the mock file path for a given file name."""
    x = mock_dir / fname
    return x


def mock_read_file(mock_dir: Path, fname: str) -> str:
    fname = get_file_mock_path(mock_dir, fname)
    data = raw_read_file(fname)
    return data
