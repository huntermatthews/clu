"""Configure pytest itself
(we just cheat and grab the rootdir)
"""

import json
from pathlib import Path

import pytest

import tests


@pytest.hookimpl
def pytest_configure(config: pytest.Config) -> None:
    pytest.mock_dir = config.rootdir / "tests" / "mock_data"  # type: ignore reportAttributeAccessIssue
    setattr(tests, "mock_data_dir", Path(config.rootdir / "tests" / "mock_data"))  # type: ignore reportAttributeAccessIssue


@pytest.fixture
def host_json_loader():
    """Loads host data from JSON file"""

    def _loader(mock_host: str) -> dict:
        host_json = tests.mock_data_dir / Path(mock_host)
        host_json = host_json.with_suffix(".json")
        with open(host_json, "r") as f:
            # print(f"Loading host JSON data from: {host_json}")
            data = json.load(f)
        return data

    return _loader
