import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.proc_cpuinfo import ProcCpuinfo

from tests import dict_subset, mock_read_file, mock_data_dir


input_keys = [
    "phy.arch",
]
output_keys = [
    "phy.arch",
    "phy.cpu.arch_version",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_proc_cpuinfo_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.proc_cpuinfo.text_file") as mrf:
        mrf.side_effect = lambda fname: mock_read_file(mock_data_dir / mock_host, fname)

        facts = Facts()
        facts.update(host_input_facts)
        proc_cpuinfo = ProcCpuinfo()
        proc_cpuinfo.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
