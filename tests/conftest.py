"""Configure pytest itself"""

import json
from pathlib import Path

import pytest

import tests


@pytest.hookimpl
def pytest_configure(config: pytest.Config) -> None:
    # FIXED: I have no idea why the .rootdir is flagged in the next line
    tests.mock_data_dir = config.rootdir / "tests" / "mock_data"  # type: ignore reportAttributeAccessIssue


@pytest.fixture
def host_json_loader():
    """Loads host data from JSON file"""

    def _loader(mock_host: str) -> dict:
        # FIXED: the Path() in the next line works around a pytest LocalPath problem.
        host_json = Path(tests.mock_data_dir / mock_host)
        host_json = host_json.with_suffix(".json")
        with open(host_json, "r") as f:
            data = json.load(f)
        return data

    return _loader
