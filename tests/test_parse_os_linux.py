import pytest
from unittest.mock import patch

from clu.os_linux import parse_os_linux

from tests import mock_read_program, mock_read_file
import tests.test_parse_clu


@pytest.mark.parametrize(
    "mock_host, expected_result",
    [
        (
            "host1",
            {
                "os.name": "Linux",
                "os.kernel.name": "Linux",
                "os.hostname": "<hostname1>",
                "os.kernel.version": "5.14.0-578.el9.x86_64",
                "os.distro.name": "centos",
                "os.distro.version": "9",
                "phy.arch": "x86_64",
                "phy.platform": "vmware",
                "phy.cpu.arch_version": "x86_64_v4",
                "phy.ram": "4.0 GB",
                "phy.cpu.model": "Intel(R) Xeon(R) Gold 6342 CPU @ 2.80GHz",
                "phy.cpu.vendor": "GenuineIntel",
                "phy.cpu.cores": "2",
                "phy.cpu.threads": "2",
                "phy.cpu.sockets": "1",
                "os.selinux.enable": "True",
                "os.selinux.mode": "Permissive",
                "salt.no_salt.exists": "True",
                "salt.no_salt.reason": "matthewsht - 2025-05-04 - This server is full of sadness.",
                "run.uptime": "1 hour, 1 second",
            },
        ),
        (
            "host2",
            {
                "os.name": "Linux",
                "os.kernel.name": "Linux",
                "os.hostname": "<hostname2>",
                "os.kernel.version": "6.4.0-150600.23.47-default",
                "os.distro.name": "opensuse-leap",
                "os.distro.version": "15.6",
                "phy.arch": "x86_64",
                "phy.platform": "vmware",
                "phy.cpu.arch_version": "x86_64_v4",
                "phy.ram": "4.0 GB",
                "phy.cpu.model": "Intel(R) Xeon(R) Gold 6342 CPU @ 2.80GHz",
                "phy.cpu.vendor": "GenuineIntel",
                "phy.cpu.cores": "2",
                "phy.cpu.threads": "2",
                "phy.cpu.sockets": "1",
                "os.selinux.enable": "Unknown/Error",
                "os.selinux.mode": "Unknown/Error",
                "salt.no_salt.exists": "False",
                "run.uptime": "28 days, 6 hours, 37 minutes, 24 seconds",
            },
        ),
        (
            "host3",
            {
                "os.name": "Linux",
                "os.kernel.name": "Linux",
                "os.hostname": "<host3>",
                "os.kernel.version": "5.14.0-503.35.1.el9_5.x86_64",
                "os.distro.name": "rhel",
                "os.distro.version": "9.5",
                "phy.arch": "x86_64",
                "phy.platform": "physical",
                "sys.vendor": "Dell Inc.",
                "sys.model.family": "PowerEdge",
                "sys.model.name": "PowerEdge R660xs",
                "sys.serial_no": "95KQ144",
                "sys.uuid": "4c4c4544-0035-4b10-8051-b9c04f313434",
                "sys.oem": "Dell Inc.",
                "sys.asset_no": "0123456789",
                "phy.cpu.arch_version": "x86_64_v4",
                "phy.ram": "64.0 GB",
                "phy.cpu.model": "Intel(R) Xeon(R) Silver 4410Y",
                "phy.cpu.vendor": "GenuineIntel",
                "phy.cpu.cores": "24",
                "phy.cpu.threads": "48",
                "phy.cpu.sockets": "2",
                "os.selinux.enable": "False",
                "os.selinux.mode": "Disabled",
                "salt.no_salt.exists": "False",
                "run.uptime": "4 days, 1 hour, 25 minutes, 35 seconds",
                "bmc.ipv4_source": "DHCP Address",
                "bmc.ipv4_address": "10.136.28.96",
                "bmc.ipv4_mask": "255.255.255.0",
                "bmc.mac_address": "d0:46:0c:70:51:b8",
                "bmc.firmware_version": "7.10",
                "bmc.manufacturer_id": "674",
                "bmc.manufacturer_name": "DELL Inc",
            },
        ),
    ],
)
def test_parse_os_linux(mock_host, expected_result):
    """Test parse_os_linux function with mock data from different hosts."""

    expected_result.update(tests.test_parse_clu.expected_result)

    # patch out parse_clu() because we test it elsewhere and its output varies too much.
    with patch("clu.os_linux.read_program") as lmrp, patch("clu.os_linux.read_file") as lmrf, patch(
        "clu.os_generic.read_program"
    ) as gmrp, patch("clu.os_linux.parse_clu") as cpc:
        lmrp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)
        lmrf.side_effect = lambda filepath: mock_read_file(pytest.mock_dir / mock_host, filepath)
        gmrp.side_effect = lambda cmdline: mock_read_program(pytest.mock_dir / mock_host, cmdline)
        cpc.side_effect = lambda facts: facts.update(tests.test_parse_clu.expected_result)

        facts = parse_os_linux()

        # Assert the expected results
        assert facts == expected_result, mock_host
