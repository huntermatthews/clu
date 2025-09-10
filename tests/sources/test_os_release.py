import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.os_release import OsRelease
from clu.sources import PARSE_FAIL_MSG

from tests import mock_read_file


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"os.distro.name": "centos", "os.distro.version": "9"}),
        ("host2", {"os.distro.name": "opensuse-leap", "os.distro.version": "15.6"}),
        ("host3", {"os.distro.name": "rhel", "os.distro.version": "9.5"}),
        ("macos", {"os.distro.name": PARSE_FAIL_MSG, "os.distro.version": PARSE_FAIL_MSG}),
    ],
)
def test_os_release_parse(mock_host, expected_result):
    """Test parse_os_release function with mock data from different hosts."""

    with patch("clu.sources.os_release.text_file") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_file(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        os_release = OsRelease()
        os_release.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
