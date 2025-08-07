# Pytest


we'd setup patch's and then:

```python
# test_cli_script.py
import subprocess

def test_cli_script_output():
    result = subprocess.run(["python", "my_cli_script.py"], capture_output=True, text=True)
    assert result.stdout.strip() == "Hello, World!"
    assert result.returncode == 0
```
