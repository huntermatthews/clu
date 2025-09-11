import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.system_version_plist import SystemVersionPlist

from tests import dict_subset, mock_read_file, mock_data_dir

input_keys = []
output_keys = [
    "os.name",
    "os.version",
    "os.build",
    "id.build_id",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("macos", input_keys, output_keys),
    ],
)
def test_system_version_plist_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.system_version_plist.text_file") as mrf:
        mrf.side_effect = lambda fname: mock_read_file(mock_data_dir / mock_host, fname)

        facts = Facts()
        facts.update(host_input_facts)
        system_version_plist = SystemVersionPlist()
        system_version_plist.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
