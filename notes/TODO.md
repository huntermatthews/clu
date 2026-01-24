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
