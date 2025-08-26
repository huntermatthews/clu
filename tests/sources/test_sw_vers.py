import pytest
from unittest.mock import patch

from clu import Facts
from clu.sources.sw_vers import SwVers

from tests import mock_read_program


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "host1",
            {
                "os.build": "Unknown/Error",
                "os.name": "Unknown/Error",
                "os.version": "Unknown/Error",
            },
        ),
        (
            "host2",
            {
                "os.build": "Unknown/Error",
                "os.name": "Unknown/Error",
                "os.version": "Unknown/Error",
            },
        ),
        (
            "host3",
            {
                "os.build": "Unknown/Error",
                "os.name": "Unknown/Error",
                "os.version": "Unknown/Error",
            },
        ),
        ("macos", {"os.name": "macOS", "os.version": "15.5", "os.build": "24F74"}),
    ],
)
def test_sw_vers_parse(mock_host, expected_result):
    """Test parse_sw_vers function with mock data from different hosts."""

    with patch("clu.sources.sw_vers.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        sw_vers = SwVers()
        sw_vers.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
