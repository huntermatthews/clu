import pytest
from unittest.mock import patch

from clu.sources.ip_addr import IpAddr
from clu.facts import Facts

from tests import dict_subset, mock_read_program, mock_data_dir


input_keys = []
output_keys = ["net.macs", "net.ipv4", "net.ipv6", "net.devs"]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_ip_addr_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.ip_addr.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(host_input_facts)
        ip_addr = IpAddr()
        ip_addr.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
