import pytest
from unittest.mock import patch

from clu.os_darwin import parse_os_darwin

from tests import mock_read_program
import tests.test_parse_clu


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "macos",
            {
                "sys.vendor": "Apple",
                "os.kernel.name": "Darwin",
                "os.hostname": "<hostname2>",
                "os.kernel.version": "24.3.0",
                "phy.arch": "arm64",
                "os.name": "macOS",
                "os.version": "15.5",
                "os.build": "24F74",
                "os.code_name": "Sequoia",
                "run.uptime": "1:03",
            },
        ),
    ],
)
def test_parse_os_darwin(mock_host, expected_result):
    """Test parse_os_darwin function with mock data from different hosts."""

    expected_result.update(tests.test_parse_clu.expected_result)

    # patch out parse_clu() because we test it elsewhere and its output varies too much.
    with patch("clu.os_darwin.read_program") as drp, patch(
        "clu.os_generic.read_program"
    ) as grp, patch("clu.os_darwin.parse_clu") as cpc:
        drp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)
        grp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)
        cpc.side_effect = lambda facts: facts.update(tests.test_parse_clu.expected_result)

        facts = parse_os_darwin()

        # Assert the expected results
        assert facts == expected_result, mock_host
