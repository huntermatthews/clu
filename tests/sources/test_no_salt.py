import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.no_salt import NoSalt

from tests import dict_subset, mock_read_file, mock_data_dir

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
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.no_salt.text_file") as mrf:
        mrf.side_effect = lambda fname, optional: mock_read_file(
            mock_data_dir / mock_host, fname, optional
        )

        facts = Facts()
        facts.update(host_input_facts)
        no_salt = NoSalt()
        no_salt.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
