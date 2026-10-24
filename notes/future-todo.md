# TODO

## Boot Session Identifiers

- **macOS**: `kern.bootsessionuuid` - Session UUID that changes every boot
- **Linux**: `/proc/sys/kernel/random/boot_id` - UUID that changes every boot
- **FreeBSD**: Boot session as epoch generated from `kern.boottime`
- **Windows**: Boot session as epoch from System Boot Time

## Machine Identifiers

- **Windows**: Read MachineGuid from registry `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography\MachineGuid` (Go has `golang.org/x/sys/windows/registry` package)

## Windows Registry Suffering

**System Identification:**
1. `HKLM\SYSTEM\CurrentControlSet\Control\ComputerName\ComputerName\ComputerName` - Full computer name
2. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\EditionID` - Windows edition (Pro, Enterprise, Home)
3. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\DisplayVersion` - Display version (22H2, 23H2)
4. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProductName` - Product name
5. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ReleaseId` - Release ID
6. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\InstallationType` - Client/Server

**Hardware/System:**
7. `HKLM\SYSTEM\CurrentControlSet\Control\SystemInformation\SystemManufacturer` - System manufacturer
8. `HKLM\SYSTEM\CurrentControlSet\Control\SystemInformation\SystemProductName` - System model
9. `HKLM\SYSTEM\CurrentControlSet\Control\SystemInformation\BIOSVersion` - BIOS version
10. `HKLM\HARDWARE\DESCRIPTION\System\CentralProcessor\0\ProcessorNameString` - CPU name
11. `HKLM\SYSTEM\CurrentControlSet\Control\TimeZoneInformation\TimeZoneKeyName` - Time zone

**Network:**
12. `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Domain` - DNS domain
13. `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Hostname` - TCP/IP hostname
14. `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\SearchList` - DNS search list

**Windows Activation/Licensing:**
15. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\RegisteredOwner` - Registered owner
16. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\RegisteredOrganization` - Registered org
17. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\DigitalProductId` - Product ID (binary)

**Updates/Patches:**
18. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\UBR` - Update Build Revision
19. `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update\AUOptions` - Auto-update setting

**System Configuration:**
20. `HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment\PROCESSOR_ARCHITECTURE` - Processor arch
21. `HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment\NUMBER_OF_PROCESSORS` - Processor count
22. `HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management\PagingFiles` - Page file config
23. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\InstallDate` - Install date (epoch)

**Virtualization:**
24. `HKLM\SYSTEM\CurrentControlSet\Services\Disk\Enum\0` - First disk (hints at VM: VBOX, VMware)
25. `HKLM\HARDWARE\DESCRIPTION\System\SystemBiosVersion` - BIOS version (may contain VM identifiers)

**User Environment:**
26. `HKLM\SYSTEM\CurrentControlSet\Control\ComputerName\ActiveComputerName\ComputerName` - Active computer name
27. `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\SystemRoot` - Windows directory

## macOS sysctl Keys

**System Identification:**
1. `kern.ostype` - OS type (Darwin)
2. `kern.osrelease` - Darwin version (23.6.0)
3. `kern.osversion` - Full kernel version string
4. `kern.version` - Complete kernel version with build info
5. `kern.hostname` - System hostname
6. `kern.bootsessionuuid` - Boot session UUID (already in Boot Session section)
7. `kern.osproductversion` - macOS version (14.5, 15.0, etc.)
8. `kern.osproductversioncompat` - Compatible product version

**Hardware - CPU:**
9. `hw.model` - Machine model (MacBookPro18,3)
10. `hw.machine` - Machine architecture (arm64, x86_64)
11. `hw.ncpu` - Total number of CPUs
12. `hw.physicalcpu` - Physical CPU cores
13. `hw.logicalcpu` - Logical CPU cores (with hyperthreading)
14. `hw.cpufrequency` - CPU frequency in Hz
15. `hw.cpufrequency_max` - Max CPU frequency
16. `machdep.cpu.brand_string` - Full CPU name/model
17. `machdep.cpu.vendor` - CPU vendor (GenuineIntel, Apple)
18. `machdep.cpu.family` - CPU family
19. `machdep.cpu.model` - CPU model number
20. `machdep.cpu.core_count` - CPU core count
21. `machdep.cpu.thread_count` - CPU thread count
22. `machdep.cpu.features` - CPU feature flags

**Hardware - Memory:**
23. `hw.memsize` - Physical memory in bytes
24. `hw.pagesize` - Memory page size
25. `vm.swapusage` - Swap usage info

**Boot/Time:**
26. `kern.boottime` - Boot time (sec/usec format)
27. `kern.clockrate` - Clock rate info

**System Configuration:**
28. `kern.maxproc` - Maximum processes
29. `kern.maxfiles` - Maximum open files
30. `kern.securelevel` - Security level

**Virtualization/Hypervisor:**
31. `kern.hv_support` - Hypervisor framework support (1 = yes, 0 = no)

**APFS/Boot (macOS-specific):**
32. `kern.bootuuid` - Boot volume UUID
33. `kern.apfsprebootuuid` - APFS preboot volume UUID
34. `kern.bootobjectspath` - Boot objects path

## Linux /proc and /sys Paths (Not Yet Parsed)

**System Identification:**
1. `/proc/sys/kernel/ostype` - OS type (Linux)
2. `/proc/sys/kernel/osrelease` - Kernel release
3. `/proc/sys/kernel/version` - Kernel version string
4. `/proc/sys/kernel/hostname` - System hostname
5. `/proc/sys/kernel/domainname` - Domain name
6. `/etc/machine-id` - Persistent machine ID
7. `/proc/sys/kernel/random/boot_id` - Boot session UUID (already in Boot Session section)

**Hardware - CPU:**
8. `/sys/devices/system/cpu/online` - Online CPUs list
9. `/sys/devices/system/cpu/possible` - Possible CPUs
10. `/sys/devices/system/cpu/present` - Present CPUs
11. `/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq` - Current CPU frequency
12. `/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq` - Max CPU frequency
13. `/sys/devices/system/cpu/cpu0/cpufreq/scaling_min_freq` - Min CPU frequency
14. `/sys/devices/system/cpu/cpu0/topology/core_id` - CPU core topology

**Hardware - Memory:**
15. `/proc/meminfo` - Memory statistics (detailed)
16. `/sys/devices/system/memory/block_size_bytes` - Memory block size
17. `/proc/swaps` - Swap devices

**Boot/Time:**
18. `/proc/stat` - Boot time (btime field)

**Virtualization:**
19. `/proc/cpuinfo` - Check hypervisor field specifically
20. `/sys/hypervisor/type` - Hypervisor type (Xen)
21. `/proc/xen/` - Xen-specific info
22. `/dev/kvm` - KVM availability (check existence)

**Network:**
23. `/sys/class/net/` - Network interfaces
24. `/proc/net/dev` - Network device statistics
25. `/etc/hostname` - Hostname file
26. `/etc/resolv.conf` - DNS configuration

**Security:**
27. `/proc/sys/kernel/dmesg_restrict` - dmesg access restriction
28. `/proc/sys/kernel/kptr_restrict` - Kernel pointer restriction
29. `/proc/sys/kernel/yama/ptrace_scope` - ptrace restrictions

**System Configuration:**
30. `/proc/sys/kernel/pid_max` - Maximum PID value
31. `/proc/sys/fs/file-max` - Maximum open files
32. `/proc/sys/kernel/threads-max` - Maximum threads
33. `/proc/cmdline` - Kernel boot parameters
34. `/proc/modules` - Loaded kernel modules
35. `/proc/sys/vm/swappiness` - Swap tendency (0-100)
36. `/proc/sys/vm/overcommit_memory` - Memory overcommit mode

**Container Detection:**
37. `/.dockerenv` - Docker container indicator (check existence)
38. `/run/systemd/container` - systemd container type
39. `/proc/1/cgroup` - Control groups (container hints)
