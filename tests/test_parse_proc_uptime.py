import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_proc_uptime

from tests import mock_read_file


@pytest.mark.parametrize("mock_host, expected_result", [
    ("host1", {"run.uptime": "3601"}),
    ("host2", {"run.uptime": "2443044.61"}),
    ("host3", {"run.uptime": "350735.47"}),
    ("macos", {}),
])
def test_parse_proc_uptime(mock_host, expected_result):
    """Test parse_proc_uptime function with mock data from different hosts."""

    with patch.object(config, 'debug', 0, create=True), patch('clu.os_linux.read_file') as mrf:
        mrf.return_value = mock_read_file(pytest.mock_dir / mock_host, '/proc/uptime')

        facts = Facts()
        parse_proc_uptime(facts)

        # # Assert the expected results
        assert facts == expected_result
