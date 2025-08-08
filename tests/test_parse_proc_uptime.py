import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_proc_uptime

from tests import mock_read_file


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"run.uptime": "1 hour, 1 second"}),
        ("host2", {"run.uptime": "28 days, 6 hours, 37 minutes, 24 seconds"}),
        ("host3", {"run.uptime": "4 days, 1 hour, 25 minutes, 35 seconds"}),
        ("macos", {}),
    ],
)
def test_parse_proc_uptime(mock_host, expected_result):
    """Test parse_proc_uptime function with mock data from different hosts."""

    with patch("clu.os_linux.read_file") as mrf:
        mrf.return_value = mock_read_file(pytest.mock_dir / mock_host, "/proc/uptime")

        facts = Facts()
        parse_proc_uptime(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
