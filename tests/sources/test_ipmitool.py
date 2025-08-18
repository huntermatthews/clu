import pytest
from unittest.mock import patch

from clu import Facts
from clu.sources.ipmitool import Ipmitool

from tests import dict_subset, mock_read_program, mock_data_dir


input_keys = ["phy.platform"]
output_keys = [
    "phy.platform",
    "bmc.ipv4_source",
    "bmc.ipv4_address",
    "bmc.ipv4_mask",
    "bmc.mac_address",
    "bmc.firmware_version",
    "bmc.manufacturer_id",
    "bmc.manufacturer_name",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_ipmitool_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.ipmitool.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(host_input_facts)
        ipmitool = Ipmitool()
        ipmitool.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
