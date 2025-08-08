import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_uname

from tests import mock_read_program


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
def test_parse_uname(mock_host, expected_result):
    """Test parse_uname function with mock data from different hosts."""

    with patch("clu.os_generic.read_program") as mrf:
        mrf.return_value = mock_read_program(pytest.mock_dir / mock_host, "uname -snrm")

        facts = Facts()
        parse_uname(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
