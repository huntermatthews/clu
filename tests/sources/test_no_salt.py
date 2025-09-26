import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.no_salt import NoSalt

from tests import dict_subset, mock_text_file, set_mock_dir

input_keys = []
output_keys = [
    "salt.no_salt.exists",
    "salt.no_salt.reason",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_no_salt_parse(mock_host, input_keys, output_keys, host_json_loader):
    set_mock_dir(mock_host)
    host_all_facts = host_json_loader()

    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.input.raw_text_file", new=mock_text_file):
        facts = Facts()
        facts.update(host_input_facts)
        no_salt = NoSalt()
        no_salt.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
