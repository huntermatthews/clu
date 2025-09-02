import pytest

from clu import Facts
from clu.sources.macos_name import MacOSName
from clu.sources import PARSE_FAIL_MSG


@pytest.mark.parametrize(
    "input_facts, expected_result",
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
def test_parse_macos_name(input_facts, expected_result):
    """Test parse_macos_name function with mock data from various versions."""

    facts = Facts()
    macos_name = MacOSName()
    facts.update(input_facts)
    macos_name.parse(facts)

    # Assert the expected results
    assert facts == expected_result, input_facts
