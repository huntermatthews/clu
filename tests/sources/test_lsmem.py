import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.lsmem import Lsmem
from clu.sources import PARSE_FAIL_MSG

from tests import mock_read_program, mock_data_dir


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"phy.ram": "4.0 GB"}),
        ("host2", {"phy.ram": "4.0 GB"}),
        ("host3", {"phy.ram": "64.0 GB"}),
        ("macos", {"phy.ram": PARSE_FAIL_MSG}),
    ],
)
def test_lsmem_parse(mock_host, expected_result):
    """Test parse_lsmem function with mock data from different hosts."""

    with patch("clu.sources.lsmem.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        lsmem = Lsmem()
        lsmem.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
