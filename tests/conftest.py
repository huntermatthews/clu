""" Configure pytest itself
(we just cheat and grab the rootdir)
"""
import pytest

@pytest.hookimpl
def pytest_configure(config: pytest.Config) -> None:
    pytest.mock_dir = config.rootdir / "mock_data"
