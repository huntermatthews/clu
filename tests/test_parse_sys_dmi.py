import pytest
from unittest.mock import patch

from clu import Facts
from clu.os_linux import parse_sys_dmi

from tests import mock_read_file


@pytest.mark.parametrize(
    "mock_host, input_facts, expected_result",
    [
        ("host1", {"phy.platform": "vmware"}, {"phy.platform": "vmware"}),
        ("host2", {"phy.platform": "vmware"}, {"phy.platform": "vmware"}),
        (
            "host3",
            {"phy.platform": "physical"},
            {
                "phy.platform": "physical",
                "sys.vendor": "Dell Inc.",
                "sys.model.family": "PowerEdge",
                "sys.model.name": "PowerEdge R660xs",
                "sys.serial_no": "95KQ144",
                "sys.uuid": "4c4c4544-0035-4b10-8051-b9c04f313434",
                "sys.oem": "Dell Inc.",
                "sys.asset_no": "0123456789",
            },
        ),
    ],
)
def test_parse_sys_dmi(mock_host, input_facts, expected_result):
    """Test parse_sys_dmi function with mock data from different hosts."""

    with patch("clu.os_linux.text_file") as mrf:
        mrf.side_effect = lambda cmdline: mock_read_file(pytest.mock_dir / mock_host, cmdline)

        facts = Facts()
        facts.update(input_facts)
        parse_sys_dmi(facts)

        # Assert the expected results
        assert facts == expected_result, mock_host
