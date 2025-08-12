import pytest
from unittest.mock import patch

from clu import Facts
from clu.os_linux import parse_ipmitool

from tests import mock_read_program


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "host1",
            {
                "phy.platform": "vmware",
            },
        ),
        (
            "host2",
            {
                "phy.platform": "vmware",
            },
        ),
        (
            "host3",
            {
                "phy.platform": "physical",
                "bmc.ipv4_source": "DHCP Address",
                "bmc.ipv4_address": "10.136.28.96",
                "bmc.ipv4_mask": "255.255.255.0",
                "bmc.mac_address": "d0:46:0c:70:51:b8",
                "bmc.firmware_version": "7.10",
                "bmc.manufacturer_id": "674",
                "bmc.manufacturer_name": "DELL Inc",
            },
        ),
        (
            "macos",
            {
                "phy.platform": "Unknown/Error",
            },
        ),
    ],
)
def test_parse_ipmitool(mock_host, expected_result):
    """Test parse_ipmitool function with mock data from different hosts."""

    with patch("clu.os_linux.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        parse_ipmitool(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
