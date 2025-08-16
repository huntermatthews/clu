import pytest
from unittest.mock import patch

from clu import Facts
from clu.sources.proc_uptime import ProcUptime

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
def test_proc_uptime_parse(mock_host, expected_result):
    """Test parse_proc_uptime function with mock data from different hosts."""

    with patch("clu.sources.proc_uptime.text_file") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_file(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        proc_uptime = ProcUptime()
        proc_uptime.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
