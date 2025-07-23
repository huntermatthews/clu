"""Doc Incomplete."""

import argparse



config = argparse.Namespace()


def recursive_dict_update(target_dict, source_dict):
    """
    Recursively updates target_dict with values from source_dict.

    Args:
        target_dict (dict): The dictionary to be updated.
        source_dict (dict): The dictionary containing the updates.
    """
    for key, value in source_dict.items():
        if key in target_dict and isinstance(target_dict[key], dict) and isinstance(value, dict):
            # If both are dictionaries, recurse
            recursive_dict_update(target_dict[key], value)
        else:
            # Otherwise, update or add the value
            target_dict[key] = value

# ------
# https://github.com/pydantic/pydantic/blob/main/pydantic/_internal/_utils.py
# SPDX-License-Identifier: MIT
from typing import Any, TypeVar

KeyType = TypeVar('KeyType')


def deep_update(mapping: dict[KeyType, Any], *updating_mappings: dict[KeyType, Any]) -> dict[KeyType, Any]:
    updated_mapping = mapping.copy()
    for updating_mapping in updating_mappings:
        for k, v in updating_mapping.items():
            if k in updated_mapping and isinstance(updated_mapping[k], dict) and isinstance(v, dict):
                updated_mapping[k] = deep_update(updated_mapping[k], v)
            else:
                updated_mapping[k] = v
    return updated_mapping
