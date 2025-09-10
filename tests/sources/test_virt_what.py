import pytest
from unittest.mock import patch

from clu.facts import Facts
from clu.sources.virt_what import VirtWhat
from clu.sources import PARSE_FAIL_MSG

from tests import mock_read_program


# TODO: virt-what needs to more carefully test rc!=0, but for now we just check the data
@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        ("host1", {"phy.platform": "vmware"}),
        ("host2", {"phy.platform": "vmware"}),
        ("host3", {"phy.platform": "physical"}),
        ("macos", {"phy.platform": PARSE_FAIL_MSG}),
    ],
)
def test_virt_what_parse(mock_host, expected_result):
    """Test parse_virt_what function with mock data from different hosts."""

    with patch("clu.sources.virt_what.text_program") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        virt_what = VirtWhat()
        virt_what.parse(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
