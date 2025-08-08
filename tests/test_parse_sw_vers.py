import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_darwin import parse_sw_vers

from tests import mock_read_program


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {}),
        ("host2", {}),
        ("host3", {}),
        ("macos", {"os.name": "macOS", "os.version": "15.5", "os.build": "24F74"}),
    ],
)
def test_parse_sw_vers(mock_host, expected_result):
    """Test parse_sw_vers function with mock data from different hosts."""

    with patch("clu.os_darwin.read_program") as mrf:
        mrf.return_value = mock_read_program(pytest.mock_dir / mock_host, "sw_vers")

        facts = Facts()
        parse_sw_vers(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
