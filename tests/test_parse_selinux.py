import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_selinux

from tests import mock_read_program


@pytest.mark.parametrize("mock_host, expected_result", [
    ("host1", {"os.selinux.enable": "True", "os.selinux.mode": "Permissive"}),
    ("host2", {"os.selinux.enable": "Unknown/Error", "os.selinux.mode": "Unknown/Error"}),
    ("host3", {"os.selinux.enable": "False", "os.selinux.mode": "Disabled"}),
    ("macos", {"os.selinux.enable": "Unknown/Error", "os.selinux.mode": "Unknown/Error"}),
])
def test_parse_selinux(mock_host, expected_result):
    """Test parse_selinux function with mock data from different hosts."""

    with patch.object(config, 'debug', 0, create=True), patch('clu.os_linux.read_program') as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        parse_selinux(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
