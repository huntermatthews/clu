import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.lscpu import Lscpu

from tests import dict_subset, mock_text_program, set_mock_dir

input_keys = []
output_keys = [
    "phy.cpu.model",
    "phy.cpu.vendor",
    "phy.cpu.cores",
    "phy.cpu.threads",
    "phy.cpu.sockets",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_lscpu_parse(mock_host, input_keys, output_keys, host_json_loader):
    set_mock_dir(mock_host)
    host_all_facts = host_json_loader()

    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.input.raw_text_program", new=mock_text_program):
        facts = Facts()
        facts.update(host_input_facts)
        lscpu = Lscpu()
        lscpu.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
