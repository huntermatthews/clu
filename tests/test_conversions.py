import pytest


from clu.conversions import bytes_to_si, si_to_bytes


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
