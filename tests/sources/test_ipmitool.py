import pytest
from unittest.mock import patch

from clu import Facts
from clu.sources.ipmitool import Ipmitool

from tests import mock_read_program


@pytest.mark.parametrize(
    "mock_host, input_facts, expected_result",
    [
        (
            "host1",
            {
                "phy.platform": "vmware",
            },
            {
                "phy.platform": "vmware",
            },
        ),
        (
            "host2",
            {
                "phy.platform": "vmware",
            },
            {
                "phy.platform": "vmware",
            },
        ),
        (
            "host3",
            {
                "phy.platform": "physical",
            },
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
            {
                "phy.platform": "Unknown/Error",
            },
        ),
    ],
)
def test_ipmitool_parse(mock_host, input_facts, expected_result):
    with patch("clu.sources.ipmitool.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(input_facts)
        ipmitool = Ipmitool()
        ipmitool.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
