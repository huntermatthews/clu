import pytest


from clu.os_linux import parse_os_release

@pytest.mark.parametrize("mock_name, expected_output",
    [
        ("host1", {'os.distro.name': 'centos',         'os.distro.version': '9'}),
        ("host2", {'os.distro.name': 'open-suse-leap', 'os.distro.version': '15.6'}),
        ("host3", {'os.distro.name': 'rhel',           'os.distro.version': '9.5'}),
    ]
)
def test_parse_os_release(mock_name, expected_output):
    call_thing(mock_name)
    assert parse_os_release() == expected_output


# @pytest.fixture
# def shared_datadir(request: pytest.FixtureRequest, tmp_path: Path) -> Path:
#     original_shared_path = os.path.join(request.fspath.dirname, "data")
#     temp_path = tmp_path / "data"
#     shutil.copytree(
#         _win32_longpath(original_shared_path), _win32_longpath(str(temp_path))
#     )
#     return temp_path


# @pytest.fixture(scope="module")
# def original_datadir(request: pytest.FixtureRequest) -> Path:
#     return Path(request.path).with_suffix("")


# @pytest.fixture
# def datadir(original_datadir: Path, tmp_path: Path) -> Path:
#     result = tmp_path / original_datadir.stem
#     if original_datadir.is_dir():
#         shutil.copytree(
#             _win32_longpath(str(original_datadir)), _win32_longpath(str(result))
#         )
#     else:
#         result.mkdir()
#     return result
