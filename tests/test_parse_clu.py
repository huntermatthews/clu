import pytest
from unittest.mock import patch
from unittest import mock

from clu.facts import Facts
from clu.os_generic import parse_clu

mock_sys = mock.MagicMock()
mock_sys.argv = ['/some/path/clu', '--test']
mock_sys.executable = '/usr/bin/pypy'
mock_sys.version_info = (3, 8, 5)
mock_sys.getcwd.return_value = '/some/working/directory'
mock_sys.__version__ = '9.9.9'

expected_result = {
    "clu.binary": '/some/path/clu',
    "clu.version": '9.9.9',
    "clu.python.binary": '/usr/bin/pypy',
    "clu.python.version": '3.8.5',
    "clu.cmdline": '/some/path/clu --test',
    "clu.cwd": '/some/working/directory',
}

def test_parse_clu():
    """Test parse_clu function.

    We can mock this from the hosts when I figure out how I want to collect this all.
    """

    with patch('sys.argv', mock_sys.argv), \
         patch('sys.executable', mock_sys.executable), \
         patch('sys.version_info', mock_sys.version_info), \
         patch('os.getcwd', mock_sys.getcwd), \
         patch('clu.os_generic.__about__', mock_sys):

        facts = Facts()
        parse_clu(facts)

        # Assert the expected results
        assert facts == expected_result
