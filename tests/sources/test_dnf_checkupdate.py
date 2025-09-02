import pytest
from unittest.mock import patch

from clu.sources.dnf_checkupdate import DnfCheckUpdate
from clu import Facts
from clu.config import set_config, Namespace

from tests import dict_subset, mock_read_program, mock_data_dir

input_keys = []
output_keys = ["run.update_required"]


def set_config_for_tests():
    args = Namespace()
    args.net = True
    set_config(args)


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_dnf_checkupdate_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    set_config_for_tests()

    with patch("clu.sources.dnf_checkupdate.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(host_input_facts)
        dnf_checkupdate = DnfCheckUpdate()
        dnf_checkupdate.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
