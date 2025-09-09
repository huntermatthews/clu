import pytest
from unittest.mock import patch

from clu import facts, Facts
from clu.sources.os_release import OsRelease
from clu.sources import PARSE_FAIL_MSG

from tests import mock_read_file, mock_data_dir


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
        print(mock_data_dir)
        mrf.side_effect = lambda fname: mock_read_file(mock_data_dir / mock_host, fname)

        expected_facts = Facts()
        expected_facts.update(expected_result)

        os_release = OsRelease()
        os_release.parse()

        # Assert the expected results
        assert isinstance(facts, Facts)
        assert isinstance(expected_facts, Facts)
        assert facts == expected_facts, mock_host
