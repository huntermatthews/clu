import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_virt_what

from tests import mock_read_program


# TODO: virt-what needs to more carefully test rc!=0, but for now we just check the data
@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"phy.platform": "vmware"}),
        ("host2", {"phy.platform": "vmware"}),
        ("host3", {"phy.platform": "physical"}),
        ("macos", {"phy.platform": "Unknown/Error"}),
    ],
)
def test_parse_virt_what(mock_host, expected_result):
    """Test parse_virt_what function with mock data from different hosts."""

    with patch("clu.os_linux.read_program") as mrf:
        mrf.return_value = mock_read_program(pytest.mock_dir / mock_host, "virt-what")

        facts = Facts()
        parse_virt_what(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
