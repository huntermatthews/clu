import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_udevadm_ram

from tests import mock_read_program


@pytest.mark.parametrize("mock_host, expected_result", [
    ("host1", {"phy.ram": "4.0 GB"}),
    ("host2", {"phy.ram": "4.0 GB"}),
    ("host3", {"phy.ram": "64.0 GB"}),
    ("macos", {}),
])
def test_parse_udevadm_ram(mock_host, expected_result):
    """Test parse_udevadm_ram function with mock data from different hosts."""

    with patch.object(config, 'debug', 0, create=True), patch('clu.os_linux.read_program') as mrf:
        mrf.return_value = mock_read_program(pytest.mock_dir / mock_host, 'udevadm info --path /devices/virtual/dmi/id')

        facts = Facts()
        parse_udevadm_ram(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
