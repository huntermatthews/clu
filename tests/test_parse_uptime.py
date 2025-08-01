import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_generic import parse_uptime


from tests import mock_read_program


@pytest.mark.parametrize("mock_host, expected_result", [
    ("host1", {"run.uptime": "58 days, 21:35"}),
    ("host2", {"run.uptime": "54 days,  9:39, "}),
    ("host3", {"run.uptime": "1 day, 5 hours, 28 minutes"}),
    ("macos", {"run.uptime": "1:03"}),
])
def test_parse_uptime(mock_host, expected_result):
    """Test parse_uptime function with mock data from different hosts."""

    with patch.object(config, 'debug', 0, create=True), patch('clu.os_generic.read_program') as mrf:
        mrf.return_value = mock_read_program(pytest.mock_dir / mock_host, 'uptime')

        facts = Facts()
        parse_uptime(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
