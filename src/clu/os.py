def os_darwin_requires():
    trace("os_darwin_requires begin")
    requires_uname()
    requires_sw_vers()
    requires_macos_name()
    requires_uptime()
    requires_gru()


def os_darwin_parse():
    trace("os_darwin_parse begin")
    ATTRS["sys.vendor"] = "Apple"
    input_uname()
    input_sw_vers()
    input_macos_name()
    input_uptime()
    input_gru()


def os_linux_requires():
    trace("os_linux_requires begin")
    requires_uname()
    requires_virt_what()
    requires_os_release()
    # BUG: only do this on physical hardware
    requires_sys_dmi()
    # BUG: only do this on x86_64/amd64 systems
    requires_cpuinfo_flags()
    requires_udevadm_ram()
    requires_lscpu()
    requires_selinux()
    requires_no_salt()
    requires_uptime()
    requires_gru()


def os_linux_parse():
    trace("os_linux_parse begin")
    ATTRS["os.name"] = "Linux"
    input_uname()
    input_virt_what()
    input_os_release()
    if ATTRS.get("phy.platform") == "physical":
        input_sys_dmi()
    if ATTRS.get("phy.arch.name") in ("x86_64", "amd64"):
        input_cpuinfo_flags()
    input_udevadm_ram()
    input_lscpu()
    input_selinux()
    input_no_salt()
    input_uptime()
    input_gru()


def os_test_requires():
    trace("os_test_requires begin")
    print("not implemented")


def os_test_parse():
    trace("os_test_parse begin")
    input_ip_addr_show()


def os_unsupported_requires():
    trace("os_unsupported_requires begin")
    panic("Unsupported OS")


def os_unsupported_parse():
    trace("os_unsupported_parse begin")
    panic("Unsupported OS")
