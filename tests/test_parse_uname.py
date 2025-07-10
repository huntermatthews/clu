import pytest

from clu.os_generic import (
    requires_uname,
    parse_uname,
)




# @pytest.mark.parametrize("uname_output, expected", [
#     ("Linux hostname 5.4.0-42-generic #46-Ubuntu SMP Fri Oct 2 10:22:24 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux", {
#         "os": "Linux",
#         "hostname": "hostname",
#         "kernel": "5.4.0-42-generic",
#         "arch": "x86_64"
#     }),
#     ("Darwin hostname 19.6.0 Darwin Kernel Version 19.6.0: Mon Aug 31 22:22:22 PDT 2020; root:xnu-6153.141.2~1/RELEASE_X86_64 x86_64", {
#         "os": "Darwin",
#         "hostname": "hostname",
#         "kernel": "19.6.0",
#         "arch": "x86_64"
#     }),
#     ("Windows hostname 10.0.19041 N/A Build 19041", {
#         "os": "Windows",
#         "hostname": "hostname",
#         "kernel": "10.0.19041",
#         "arch": "x86_64"
#     }),
# ])
# def test_parse_uname(uname_output, expected):
#     assert parse_uname(uname_output) == expected
