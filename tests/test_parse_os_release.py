import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_os_release


@pytest.mark.parametrize("mock_host, expected_distro_name, expected_distro_version", [
    ("host1", "centos", "9"),
    ("host2", "opensuse-leap", "15.6"),
    ("host3", "rhel", "9.5"),
])
def test_parse_os_release(mock_host, expected_distro_name, expected_distro_version):
    """Test parse_os_release function with mock data from different hosts."""
    # Set up mock directory path
    mock_dir = pytest.mock_dir / mock_host

    # Mock the config attributes needed by the code
    with patch.object(config, 'mock', mock_dir, create=True), \
         patch.object(config, 'debug', 0, create=True):
        # Create a Facts object to store the parsed results
        facts = Facts()

        # Call the function under test
        parse_os_release(facts)

        # Assert the expected results
        assert facts["os.distro.name"] == expected_distro_name
        assert facts["os.distro.version"] == expected_distro_version


def test_parse_os_release_missing_file():
    """Test parse_os_release function when os-release file is missing."""
    # Use a mock directory that doesn't have an os-release file
    mock_dir = pytest.mock_dir / "macos"  # macos directory doesn't have /etc/os-release

    # Mock the config attributes needed by the code
    with patch.object(config, 'mock', mock_dir, create=True), \
         patch.object(config, 'debug', 0, create=True):
        # Create a Facts object to store the parsed results
        facts = Facts()

        # Call the function under test
        parse_os_release(facts)

        # Assert the expected error results
        assert facts["os.distro.name"] == "Unknown/Error"
        assert facts["os.distro.version"] == "Unknown/Error"
