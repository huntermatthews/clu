import pytest
@pytest.hookimpl
def pytest_configure(config):
    pytest.rootdir = config.rootdir
    pytest.datadir = config.rootdir / "mock_data"

