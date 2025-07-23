
import pytest

@pytest.mark.parametrize(
    'mapping, updating_mapping, expected_mapping, msg',
    [
        (
            {'key': {'inner_key': 0}},
            {'other_key': 1},
            {'key': {'inner_key': 0}, 'other_key': 1},
            'extra keys are inserted',
        ),
        (
            {'key': {'inner_key': 0}, 'other_key': 1},
            {'key': [1, 2, 3]},
            {'key': [1, 2, 3], 'other_key': 1},
            'values that can not be merged are updated',
        ),
        (
            {'key': {'inner_key': 0}},
            {'key': {'other_key': 1}},
            {'key': {'inner_key': 0, 'other_key': 1}},
            'values that have corresponding keys are merged',
        ),
        (
            {'key': {'inner_key': {'deep_key': 0}}},
            {'key': {'inner_key': {'other_deep_key': 1}}},
            {'key': {'inner_key': {'deep_key': 0, 'other_deep_key': 1}}},
            'deeply nested values that have corresponding keys are merged',
        ),
    ],
)
def test_deep_update(mapping, updating_mapping, expected_mapping, msg):
    assert deep_update(mapping, updating_mapping) == expected_mapping, msg
