"""Doc Incomplete."""


def bytes_to_si(size: int) -> str:
    """
    Convert a number of bytes to a human-readable SI string (e.g., 1536 -> '1.5 KB').
    """
    units = ["B", "KB", "MB", "GB", "TB", "PB", "EB"]
    size = float(size)
    for unit in units:
        if size < 1024:
            return f"{size:.1f} {unit}"
        size /= 1024
    return f"{size:.1f} {units[-1]}"


def si_to_bytes(size_str: str) -> int:
    """
    Convert a human-readable SI string (e.g., '1.5 KB') to bytes (int).
    """
    units = {"B": 0, "KB": 1, "MB": 2, "GB": 3, "TB": 4, "PB": 5, "EB": 6}
    size_str = size_str.strip().upper()
    number = ""
    unit = ""
    for char in size_str:
        if char.isdigit() or char == "." or char == ",":
            number += char
        elif char.isalpha():
            unit += char
    number = float(number.replace(",", ""))
    unit = unit.strip()
    if unit not in units:
        raise ValueError(f"Unknown unit: {unit}")
    exponent = units[unit]
    return int(number * (1024**exponent))


# TODO:
# - Consider adding support for binary prefixes (KiB, MiB, etc.) if needed
# - extend to ZB, YB, RB, QB -- and link to https://en.wikipedia.org/wiki/Binary_prefix
