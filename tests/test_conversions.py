import pytest


from clu.conversions import bytes_to_si, si_to_bytes


def test_bytes_to_si():
    assert bytes_to_si(1536) == "1.5 KB"
    assert bytes_to_si(1048576) == "1.0 MB"
    assert bytes_to_si(1073741824) == "1.0 GB"
    assert bytes_to_si(1099511627776) == "1.0 TB"
    assert bytes_to_si(1125899906842624) == "1.0 PB"
    assert bytes_to_si(1152921504606846976) == "1.0 EB"
    assert bytes_to_si(0) == "0.0 B"
    assert bytes_to_si(1023) == "1023.0 B"


def test_si_to_bytes():
    assert si_to_bytes("0.0 B") == 0
    assert si_to_bytes("0 B") == 0
    assert si_to_bytes("1023 B") == 1023

    assert si_to_bytes("1.5 KB") == 1536
    assert si_to_bytes("1.0 MB") == 1048576
    assert si_to_bytes("1.0 GB") == 1073741824
    assert si_to_bytes("1.0 TB") == 1099511627776
    assert si_to_bytes("1.0 PB") == 1125899906842624
    assert si_to_bytes("1.0 EB") == 1152921504606846976

    # assert si_to_bytes("1 KB") == 1000
    # assert si_to_bytes("1 MB") == 1000000
    # assert si_to_bytes("1 GB") == 1000000000
    # assert si_to_bytes("1 TB") == 1000000000000
    # assert si_to_bytes("1 PB") == 1000000000000000
    # assert si_to_bytes("1 EB") == 1000000000000000000
