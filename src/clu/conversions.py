"""Doc Incomplete."""


def bytes_to_si(size: float) -> str:
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


def si_to_bytes(size_str: str) -> float:
    """
    Convert a human-readable SI string (e.g., '1.5 KB') to bytes (int).
    """
    units = {
        "B": 0,
        "K": 1,
        "KB": 1,
        "M": 2,
        "MB": 2,
        "G": 3,
        "GB": 3,
        "T": 4,
        "TB": 4,
        "P": 5,
        "PB": 5,
        "E": 6,
        "EB": 6,
    }
    size_str = size_str.strip().upper()
    number_str = ""
    unit = ""
    for char in size_str:
        if char.isdigit() or char == "." or char == ",":
            number_str += char
        elif char.isalpha():
            unit += char
    number = float(number_str.replace(",", ""))
    unit = unit.strip()
    if unit not in units:
        raise ValueError(f"Unknown unit: {unit}")
    exponent = units[unit]
    return float(number * (1024**exponent))


def seconds_to_text(secs: int) -> str:
    """
    Convert seconds to a human-readable string with weeks, days, hours, and minutes.

    Args:
        secs: Number of seconds to convert

    Returns:
        A formatted string like "1 week, 2 days, 3 hours, 45 minutes"
    """
    if secs < 0:
        return "0 seconds"

    # Calculate time units
    months = secs // 2592000  # 30 * 24 * 60 * 60
    remaining = secs % 2592000

    days = remaining // 86400  # 24 * 60 * 60
    remaining = remaining % 86400

    hours = remaining // 3600  # 60 * 60
    remaining = remaining % 3600

    minutes = remaining // 60
    seconds = remaining % 60

    # Build result parts
    parts = []
    if months > 0:
        parts.append(f"{months} month{'s' if months != 1 else ''}")
    if days > 0:
        parts.append(f"{days} day{'s' if days != 1 else ''}")
    if hours > 0:
        parts.append(f"{hours} hour{'s' if hours != 1 else ''}")
    if minutes > 0:
        parts.append(f"{minutes} minute{'s' if minutes != 1 else ''}")
    if seconds > 0 or len(parts) == 0:
        parts.append(f"{seconds} second{'s' if seconds != 1 else ''}")

    return ", ".join(parts)
