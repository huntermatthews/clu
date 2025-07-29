import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_lscpu

from tests import mock_read_program


@pytest.mark.parametrize("mock_host, expected_result", [
    ("host1", {"phy.cpu.model": "Intel(R) Xeon(R) Gold 6342 CPU @ 2.80GHz", "phy.cpu.vendor": "GenuineIntel", "phy.cpu.cores": "2", "phy.cpu.threads": "2", "phy.cpu.sockets": "1"}),
    ("host2", {"phy.cpu.model": "Intel(R) Xeon(R) Gold 6342 CPU @ 2.80GHz", "phy.cpu.vendor": "GenuineIntel", "phy.cpu.cores": "2", "phy.cpu.threads": "2", "phy.cpu.sockets": "1"}),
    ("host3", {"phy.cpu.model": "Intel(R) Xeon(R) Silver 4410Y", "phy.cpu.vendor": "GenuineIntel", "phy.cpu.cores": "24", "phy.cpu.threads": "48", "phy.cpu.sockets": "2"}),
    ("macos", {}),
])
def test_parse_lscpu(mock_host, expected_result):
    """Test parse_lscpu function with mock data from different hosts."""

    with patch.object(config, 'debug', 0, create=True), patch('clu.os_linux.read_program') as mrf:
        mrf.return_value = mock_read_program(pytest.mock_dir / mock_host, 'lscpu')

        facts = Facts()
        parse_lscpu(facts)

        # # Assert the expected results
        assert facts == expected_result
