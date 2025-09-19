import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.os_release import OsRelease

from tests import dict_subset, mock_read_file, mock_data_dir

input_keys = []
output_keys = [
    "os.distro.name",
    "os.distro.version",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("host1", input_keys, output_keys),
        ("host2", input_keys, output_keys),
        ("host3", input_keys, output_keys),
    ],
)
def test_os_release_parse(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    with patch("clu.sources.os_release.text_file") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_file(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(host_input_facts)
        os_release = OsRelease()
        os_release.parse(facts)

        # Assert the expected results
        assert facts == host_output_facts, mock_host
