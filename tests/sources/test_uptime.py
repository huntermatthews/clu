import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.uptime import Uptime

from tests import mock_read_program, mock_data_dir


# CAUTION: linux hosts use /proc/uptime which was captured DIFFERENTLY than these
# values. You can't update this to the host_json pattern...
@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"run.uptime": "58 days, 21:35"}),
        ("host2", {"run.uptime": "54 days,  9:39, "}),
        ("host3", {"run.uptime": "1 day, 5 hours, 28 minutes"}),
        ("macos", {"run.uptime": "1:03"}),
    ],
)
def test_uptime_parse(mock_host, expected_result):
    """Test parse_uptime function with mock data from different hosts."""

    with patch("clu.sources.uptime.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        uptime = Uptime()
        uptime.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
