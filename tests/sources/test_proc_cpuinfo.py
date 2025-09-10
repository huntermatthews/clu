import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.proc_cpuinfo import ProcCpuinfo

from tests import mock_read_file


@pytest.mark.parametrize(
    "mock_host, input_facts, expected_result",
    [
        (
            "host1",
            {"phy.arch": "x86_64"},
            {"phy.arch": "x86_64", "phy.cpu.arch_version": "x86_64_v4"},
        ),
        (
            "host2",
            {"phy.arch": "x86_64"},
            {"phy.arch": "x86_64", "phy.cpu.arch_version": "x86_64_v4"},
        ),
        (
            "host3",
            {"phy.arch": "x86_64"},
            {"phy.arch": "x86_64", "phy.cpu.arch_version": "x86_64_v4"},
        ),
        ("macos", {"phy.arch": "arm64"}, {"phy.arch": "arm64"}),
    ],
)
def test_proc_cpuinfo_parse(mock_host, input_facts, expected_result):
    """Test parse_cpuinfo_flags function with mock data from different hosts."""

    with patch("clu.sources.proc_cpuinfo.text_file") as mrf:
        mrf.side_effect = lambda fname: mock_read_file(pytest.mock_dir / mock_host, fname)

        facts = Facts()
        facts.update(input_facts)
        proc_cpuinfo = ProcCpuinfo()
        proc_cpuinfo.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
