import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.udevadm_ram import UdevadmRam
from clu.sources import PARSE_FAIL_MSG

from tests import mock_read_program


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"phy.ram": "4.0 GB"}),
        ("host2", {"phy.ram": "4.0 GB"}),
        ("host3", {"phy.ram": "64.0 GB"}),
        ("macos", {"phy.ram": PARSE_FAIL_MSG}),
    ],
)
def test_udevadm_ram_parse(mock_host, expected_result):
    """Test parse_udevadm_ram function with mock data from different hosts."""

    with patch("clu.sources.udevadm_ram.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        udevadm_ram = UdevadmRam()
        udevadm_ram.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
