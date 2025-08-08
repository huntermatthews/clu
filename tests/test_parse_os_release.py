import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_os_release

from tests import mock_read_file


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"os.distro.name": "centos", "os.distro.version": "9"}),
        ("host2", {"os.distro.name": "opensuse-leap", "os.distro.version": "15.6"}),
        ("host3", {"os.distro.name": "rhel", "os.distro.version": "9.5"}),
        ("macos", {"os.distro.name": "Unknown/Error", "os.distro.version": "Unknown/Error"}),
    ],
)
def test_parse_os_release(mock_host, expected_result):
    """Test parse_os_release function with mock data from different hosts."""

    with patch("clu.os_linux.read_file") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_file(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        parse_os_release(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
