import pytest
from unittest.mock import patch

from clu import Facts, facts
from clu.sources.lsmem import Lsmem

from tests import dict_subset, mock_read_program, mock_data_dir

input_keys = []
output_keys = ["phy.ram"]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_lsmem_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.lsmem.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        expected_facts = Facts()
        expected_facts.update(host_input_facts)
        expected_facts.update(host_output_facts)

        lsmem = Lsmem()
        lsmem.parse()

        # Assert the expected results
        assert isinstance(facts, Facts)
        assert isinstance(expected_facts, Facts)
        assert facts == expected_facts, mock_host
