import pytest
from unittest.mock import patch

from clu import Facts
from clu.os_linux import parse_no_salt

from tests import mock_read_file


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "host1",
            {
                "salt.no_salt.exists": "True",
                "salt.no_salt.reason": "matthewsht - 2025-05-04 - This server is full of sadness.",
            },
        ),
        ("host2", {"salt.no_salt.exists": "False"}),
        ("host3", {"salt.no_salt.exists": "False"}),
        ("macos", {"salt.no_salt.exists": "False"}),
    ],
)
def test_parse_no_salt(mock_host, expected_result):
    """Test parse_no_salt function with mock data from different hosts."""

    with patch("clu.os_linux.text_file") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_file(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        parse_no_salt(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
