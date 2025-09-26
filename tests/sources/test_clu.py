from unittest.mock import patch

from clu.facts import Facts
from clu.sources.clu import Clu

from tests import mock_clu_metadata


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
    with patch("clu.sources.clu.raw_clu_metadata", new=mock_clu_metadata):
        clu = Clu()
        facts = Facts()
        clu.parse(facts)

        expected_facts = Facts()
        expected_facts.update(expected_result)

        # Assert the expected results
        assert facts == expected_facts
