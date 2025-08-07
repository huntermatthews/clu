import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_darwin import parse_macos_name


@pytest.mark.parametrize(
    "input_facts, expected_result",
    [
        ({"os.version": "26"}, {"os.version": "26", "os.code_name": "Tahoe"}),
        ({"os.version": "25"}, {"os.version": "25", "os.code_name": "Unknown-25"}),
        ({"os.version": "16"}, {"os.version": "16", "os.code_name": "Unknown-16"}),
        ({"os.version": "15.5"}, {"os.version": "15.5", "os.code_name": "Sequoia"}),
        ({"os.version": "14.0"}, {"os.version": "14.0", "os.code_name": "Sonoma"}),
        ({"os.version": "13.0"}, {"os.version": "13.0", "os.code_name": "Ventura"}),
        ({"os.version": "12.0"}, {"os.version": "12.0", "os.code_name": "Monterey"}),
        ({"os.version": "11.0"}, {"os.version": "11.0", "os.code_name": "Big Sur"}),
        ({"os.version": "10.15"}, {"os.version": "10.15", "os.code_name": "Unknown-10"}),
        ({"os.version": "unknown"}, {"os.version": "unknown", "os.code_name": "Unknown-unknown"}),
    ],
)
def test_parse_macos_name(input_facts, expected_result):
    """Test parse_macos_name function with mock data from various versions."""

    with patch.object(config, "debug", 0, create=True):
        facts = Facts()
        facts.update(input_facts)
        parse_macos_name(facts)

        # Assert the expected results
        assert facts == expected_result, input_facts
