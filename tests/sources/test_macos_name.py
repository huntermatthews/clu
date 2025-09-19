import pytest

from clu.facts import Facts
from clu.sources.macos_name import MacOSName
from clu.sources import PARSE_FAIL_MSG

from tests import dict_subset

input_keys = ["os.version"]
output_keys = [
    "os.version",
    "os.code_name",
]


@pytest.mark.parametrize(
    "mock_host, input_keys, output_keys",
    [
        ("macos", input_keys, output_keys),
    ],
)
def test_parse_macos_name(mock_host, input_keys, output_keys, host_json_loader):
    host_all_facts = host_json_loader(mock_host)
    host_input_facts = dict_subset(host_all_facts, input_keys)
    host_output_facts = dict_subset(host_all_facts, output_keys)

    # MacOSName does not read any files/programs, so no need to mock anything here
    facts = Facts()
    facts.update(host_input_facts)
    macos_name = MacOSName()
    macos_name.parse(facts)

    # Assert the expected results
    assert facts == host_output_facts, mock_host


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
def test_parse_macos_name_additional_inputs(input_facts, expected_result):
    """Test parse_macos_name function with mock data from various versions."""

    facts = Facts()
    macos_name = MacOSName()
    facts.update(input_facts)
    macos_name.parse(facts)

    # Assert the expected results
    assert facts == expected_result, input_facts
