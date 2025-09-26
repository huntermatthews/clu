import pytest
from unittest.mock import patch

from clu.sources.dnf_checkupdate import DnfCheckUpdate
from clu.facts import Facts
from clu.config import set_config, Namespace

from tests import dict_subset, mock_text_program, set_mock_dir

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
    set_mock_dir(mock_host)
    host_all_facts = host_json_loader()

    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    set_config_for_tests()

    with patch("clu.input.raw_text_program", new=mock_text_program):
        facts = Facts()
        facts.update(host_input_facts)
        dnf_checkupdate = DnfCheckUpdate()
        dnf_checkupdate.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
