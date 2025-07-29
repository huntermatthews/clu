import pytest
from unittest.mock import patch

from clu import config
from clu.os_darwin import parse_os_darwin

from tests import mock_read_program


@pytest.mark.parametrize("mock_host, expected_result", [
    ("macos", {
        "sys.vendor": "Apple",
        "os.kernel.name": "Darwin",
        "os.hostname": "<hostname2>",
        "os.kernel.version": "24.3.0",
        "phy.arch": "arm64",
        "os.name": "macOS",
        "os.version": "15.5",
        "os.build": "24F74",
        "os.code_name": "Sequoia",
        "run.uptime": "1:03"
    }),
])
def test_parse_os_darwin(mock_host, expected_result):
    """Test parse_os_darwin function with mock data from different hosts."""

    with patch.object(config, 'debug', 0, create=True), \
         patch('clu.os_darwin.read_program') as lmrp, \
         patch('clu.os_generic.read_program') as gmrp:
        lmrp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)
        gmrp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = parse_os_darwin()

        # Remove all entries where the key starts with "clu."
        # TODO: figure out how to mock parse_clu to avoid this
        keys_to_remove = [key for key in facts.keys() if key.startswith("clu.")]
        for key in keys_to_remove:
            del facts[key]

        # Assert the expected results
        assert facts == expected_result, mock_host
