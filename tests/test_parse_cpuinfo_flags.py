import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_cpuinfo_flags

from tests import mock_read_file


@pytest.mark.parametrize("mock_host, input_facts, expected_result", [
    ("host1", {"phy.arch": "x86_64"}, {"phy.arch": "x86_64", "phy.cpu.arch_version": "x86_64_v4"}),
    ("host2", {"phy.arch": "x86_64"}, {"phy.arch": "x86_64", "phy.cpu.arch_version": "x86_64_v4"}),
    ("host3", {"phy.arch": "x86_64"}, {"phy.arch": "x86_64", "phy.cpu.arch_version": "x86_64_v4"}),
    ("macos", {"phy.arch": "arm64"}, {"phy.arch": "arm64"}),
])
def test_parse_cpuinfo_flags(mock_host, input_facts, expected_result):
    """Test parse_cpuinfo_flags function with mock data from different hosts."""

    def result_mocker(fname):
        return mock_read_file(pytest.mock_dir / mock_host, fname)

    with patch.object(config, 'debug', 0, create=True), patch('clu.os_linux.read_file') as mrf:
        # mrf.side_effect = lambda fname: mock_read_file(pytest.mock_dir / mock_host, fname)
        mrf.side_effect = result_mocker

        facts = Facts()
        facts.update(input_facts)
        parse_cpuinfo_flags(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
