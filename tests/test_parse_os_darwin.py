import pytest
from unittest.mock import patch

from clu import config
from clu.os_darwin import parse_os_darwin

from tests import mock_read_program
import tests.test_parse_clu

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

    expected_result.update(tests.test_parse_clu.expected_result)

    with patch.object(config, 'debug', 0, create=True), \
         patch('clu.os_darwin.read_program') as lmrp, \
         patch('clu.os_generic.read_program') as gmrp, \
         patch('sys.argv', tests.test_parse_clu.mock_sys.argv), \
         patch('sys.executable', tests.test_parse_clu.mock_sys.executable), \
         patch('sys.version_info', tests.test_parse_clu.mock_sys.version_info), \
         patch('os.getcwd', tests.test_parse_clu.mock_sys.getcwd), \
         patch('clu.os_generic.__about__', tests.test_parse_clu.mock_sys):
        lmrp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)
        gmrp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = parse_os_darwin()

        # Assert the expected results
        assert facts == expected_result, mock_host
