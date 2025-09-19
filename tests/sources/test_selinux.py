import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.selinux import Selinux

from tests import dict_subset, mock_read_program, mock_data_dir


input_keys = []
output_keys = [
    "os.selinux.enable",
    "os.selinux.mode",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_selinux_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.selinux.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(host_input_facts)
        selinux = Selinux()
        selinux.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
