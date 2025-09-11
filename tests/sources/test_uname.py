import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.uname import Uname

from tests import mock_read_program, mock_data_dir


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "host1",
            {
                "os.kernel.name": "Linux",
                "os.kernel.version": "5.14.0-578.el9.x86_64",
                "os.hostname": "<hostname1>",
                "phy.arch": "x86_64",
            },
        ),
        (
            "host2",
            {
                "os.kernel.name": "Linux",
                "os.kernel.version": "6.4.0-150600.23.47-default",
                "os.hostname": "<hostname2>",
                "phy.arch": "x86_64",
            },
        ),
        (
            "host3",
            {
                "os.kernel.name": "Linux",
                "os.kernel.version": "5.14.0-503.35.1.el9_5.x86_64",
                "os.hostname": "<host3>",
                "phy.arch": "x86_64",
            },
        ),
        (
            "macos",
            {
                "os.kernel.name": "Darwin",
                "os.kernel.version": "24.3.0",
                "os.hostname": "<hostname2>",
                "phy.arch": "arm64",
            },
        ),
    ],
)
def test_uname_parse(mock_host, expected_result):
    """Test parse_uname function with mock data from different hosts."""

    with patch("clu.sources.uname.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(mock_data_dir / mock_host, cmdline)

        facts = Facts()
        uname = Uname()
        uname.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
