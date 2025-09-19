import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.uname import Uname

from tests import dict_subset, mock_read_program, mock_data_dir


input_keys = []
output_keys = ["os.kernel.name", "os.kernel.version", "os.hostname", "phy.arch"]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_uname_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.uname.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(host_input_facts)
        uname = Uname()
        uname.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
