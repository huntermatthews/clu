class Facts(dict):
    """A simple facts container that behaves like a dict."""
    def __init__(self):
        self._priority = {"primary": [], "secondary": [], "other": []}
        self._facts = {}

    def add(self, priority: str, key: str, value: str) -> None:
        if priority not in self._priority.keys():
            raise ValueError(f"Invalid priority: {priority}")
        self._priority[priority].append(key)
        self._facts[key] = value

_facts: Facts = Facts()


