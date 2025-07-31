---
from unittest.mock import patch

def test_parse_args():
    testargs = ["prog", "-f", "/home/fenton/project/setup.py"]
    with patch.object(sys, 'argv', testargs):
        setup = get_setup_file()
        assert setup == "/home/fenton/project/setup.py"
---
mocked_object.assert_called_once_with(*args)
mock_pyautogui.assert_called_once_with()
