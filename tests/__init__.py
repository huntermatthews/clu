"""Common code for all unit tests."""

from pathlib import Path

from clu.readers import raw_read_file


def get_file_mock_path(mock_dir, fname):
    """Get the mock file path for a given file name."""
    x = mock_dir / fname
    return x


def mock_read_file(mock_dir: Path, fname: str) -> str:
    fname = get_file_mock_path(mock_dir, fname)
    data = raw_read_file(fname)
    return data
