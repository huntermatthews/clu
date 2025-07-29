import pytest
from unittest.mock import patch

from clu import config
from clu.facts import Facts
from clu.os_linux import parse_os_release

from tests import mock_read_file


# we don't currently have any mock data for clu runtime, so this test is commented out
