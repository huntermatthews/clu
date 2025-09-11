import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.no_salt import NoSalt

from tests import mock_read_file, mock_data_dir


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
def test_no_salt_parse(mock_host, expected_result):
    """Test parse_no_salt function with mock data from different hosts."""

    with patch("clu.sources.no_salt.text_file") as mrf:
        mrf.side_effect = lambda fname, optional: mock_read_file(
            mock_data_dir / mock_host, fname, optional
        )

        facts = Facts()
        no_salt = NoSalt()
        no_salt.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
