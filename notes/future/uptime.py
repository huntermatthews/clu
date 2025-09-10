import re

def parse_uptime_output(uptime_output):
    """
    Parses the output of the `uptime` command and returns the uptime in seconds.
    Example uptime outputs:
      ' 14:23  up 1 day,  2:34,  3 users,  load average: 0.00, 0.01, 0.05'
      ' 14:23  up 2:34,  3 users,  load average: 0.00, 0.01, 0.05'
      ' 14:23  up 3 mins,  3 users,  load average: 0.00, 0.01, 0.05'
    """
    # Find the 'up' part
    match = re.search(r'up (.*?),\s+\d+ user', uptime_output)
    if not match:
        match = re.search(r'up (.*?),\s+load average', uptime_output)
    if not match:
        return None

    uptime_str = match.group(1).strip()

    days = 0
    hours = 0
    minutes = 0

    # Check for days
    day_match = re.search(r'(\d+)\s+day[s]?', uptime_str)
    if day_match:
        days = int(day_match.group(1))
        uptime_str = uptime_str.replace(day_match.group(0), '').strip(', ')

    # Check for hours and minutes in format H:MM
    hm_match = re.search(r'(\d+):(\d+)', uptime_str)
    if hm_match:
        hours = int(hm_match.group(1))
        minutes = int(hm_match.group(2))
    else:
        # Check for minutes only
        min_match = re.search(r'(\d+)\s+min', uptime_str)
        if min_match:
            minutes = int(min_match.group(1))

    total_seconds = days * 86400 + hours * 3600 + minutes * 60
    return total_seconds
