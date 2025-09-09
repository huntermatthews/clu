import pytest

from clu import facts, Facts
from clu.sources.macos_name import MacOSName
from clu.sources import PARSE_FAIL_MSG


@pytest.mark.parametrize(
    "host_input_facts, host_output_facts",
    [
        ({"os.version": "26"}, {"os.version": "26", "os.code_name": "Tahoe"}),
        ({"os.version": "25"}, {"os.version": "25", "os.code_name": PARSE_FAIL_MSG}),
        ({"os.version": "16"}, {"os.version": "16", "os.code_name": PARSE_FAIL_MSG}),
        ({"os.version": "15.5"}, {"os.version": "15.5", "os.code_name": "Sequoia"}),
        ({"os.version": "14.0"}, {"os.version": "14.0", "os.code_name": "Sonoma"}),
        ({"os.version": "13.0"}, {"os.version": "13.0", "os.code_name": "Ventura"}),
        ({"os.version": "12.0"}, {"os.version": "12.0", "os.code_name": "Monterey"}),
        ({"os.version": "11.0"}, {"os.version": "11.0", "os.code_name": "Big Sur"}),
        ({"os.version": "10.15"}, {"os.version": "10.15", "os.code_name": PARSE_FAIL_MSG}),
        ({"os.version": "unknown"}, {"os.version": "unknown", "os.code_name": PARSE_FAIL_MSG}),
    ],
)
def test_parse_macos_name(host_input_facts, host_output_facts):
    """Test parse_macos_name function with mock data from various versions."""

    expected_facts = Facts()
    expected_facts.update(host_input_facts)
    expected_facts.update(host_output_facts)

    facts.update(host_input_facts)
    macos_name = MacOSName()
    macos_name.parse()

    # Assert the expected results
    assert isinstance(facts, Facts)
    assert isinstance(expected_facts, Facts)
    assert facts == expected_facts
