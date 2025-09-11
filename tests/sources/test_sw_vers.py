import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.sw_vers import SwVers
from clu.sources import PARSE_FAIL_MSG

from tests import mock_read_program, mock_data_dir


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "host1",
            {
                "os.build": PARSE_FAIL_MSG,
                "os.name": PARSE_FAIL_MSG,
                "os.version": PARSE_FAIL_MSG,
            },
        ),
        (
            "host2",
            {
                "os.build": PARSE_FAIL_MSG,
                "os.name": PARSE_FAIL_MSG,
                "os.version": PARSE_FAIL_MSG,
            },
        ),
        (
            "host3",
            {
                "os.build": PARSE_FAIL_MSG,
                "os.name": PARSE_FAIL_MSG,
                "os.version": PARSE_FAIL_MSG,
            },
        ),
        ("macos", {"os.name": "macOS", "os.version": "15.5", "os.build": "24F74"}),
    ],
)
def test_sw_vers_parse(mock_host, expected_result):
    """Test parse_sw_vers function with mock data from different hosts."""

    with patch("clu.sources.sw_vers.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        sw_vers = SwVers()
        sw_vers.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
