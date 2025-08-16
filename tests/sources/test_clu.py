from unittest.mock import patch
from unittest import mock

from clu import Facts
from clu.sources.clu import Clu

mock_sys = mock.MagicMock()
mock_sys.argv = ["/some/path/clu", "--test"]
mock_sys.executable = "/usr/bin/pypy"
mock_sys.version_info = (3, 8, 5)
mock_sys.getcwd.return_value = "/some/working/directory"
mock_sys.__version__ = "9.9.9"
mock_sys.getenv.return_value = "zorro"
mock_sys.now = "2023-10-01T12:00:00+00:00"

expected_result = {
    "clu.binary": "/some/path/clu",
    "clu.version": "9.9.9",
    "clu.python.binary": "/usr/bin/pypy",
    "clu.python.version": "3.8.5",
    "clu.cmdline": "/some/path/clu --test",
    "clu.cwd": "/some/working/directory",
    "clu.user": "zorro",
    "clu.date": "2023-10-01T12:00:00+00:00",
}


def test_clu_parse():
    with patch("sys.argv", mock_sys.argv), patch("sys.executable", mock_sys.executable), patch(
        "sys.version_info", mock_sys.version_info
    ), patch("os.getcwd", mock_sys.getcwd), patch("clu.sources.clu.__about__", mock_sys), patch(
        "getpass.getuser", mock_sys.getenv
    ), patch("clu.sources.clu.Clu._get_rfc3339_timestamp", return_value=mock_sys.now):
        clu = Clu()
        facts = Facts()
        clu.parse(facts)

        # Assert the expected results
        assert facts == expected_result
