import pytest


from clu.conversions import bytes_to_si, si_to_bytes, seconds_to_text


@pytest.mark.parametrize(
    "si_string, expected_output",
    [
        ("1.5 KB", 1536),
        ("1.0 MB", 1048576),
        ("1.0 GB", 1073741824),
        ("1.0 TB", 1099511627776),
        ("1.0 PB", 1125899906842624),
        ("1.0 EB", 1152921504606846976),
        ("0.0 B", 0),
        ("0 B", 0),
        ("1023 B", 1023),
    ],
)
def test_si_to_bytes(si_string, expected_output):
    bytes_count = si_to_bytes(si_string)

    assert bytes_count == expected_output, (
        f"Expected {expected_output} for {si_string}, got {bytes_count}"
    )


@pytest.mark.parametrize(
    "bytes_count, expected_output",
    [
        (1536, "1.5 KB"),
        (1048576, "1.0 MB"),
        (1073741824, "1.0 GB"),
        (1099511627776, "1.0 TB"),
        (1125899906842624, "1.0 PB"),
        (1152921504606846976, "1.0 EB"),
        (0, "0.0 B"),
        (1023, "1023.0 B"),
    ],
)
def test_bytes_to_si(bytes_count, expected_output):
    si_string = bytes_to_si(bytes_count)

    assert si_string == expected_output, (
        f"Expected {expected_output} for {bytes_count}, got {si_string}"
    )


@pytest.mark.parametrize(
    "seconds, expected_output",
    [
        (3601, "1 hour, 1 second"),
        (2443044.61, "28 days, 6 hours, 37 minutes, 24 seconds"),
        (350735.47, "4 days, 1 hour, 25 minutes, 35 seconds"),
        (5114048.68, "1 month, 29 days, 4 hours, 34 minutes, 8 seconds"),
    ],
)
def test_seconds_to_text(seconds, expected_output):
    text = seconds_to_text(int(seconds))

    assert text == expected_output, f"Expected {expected_output} for {seconds}, got {text}"
