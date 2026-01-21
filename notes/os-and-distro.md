
# OS / Distro Support Notes

## MacOS

### Uptime

Parsing uptime(4) has a regex that is unreliable and goofy - plus parsing the result.
`sysctl kern.boottime` is far more consistent.

### Ram

❯ sysctl hw.memsize
hw.memsize: 34359738368

### SystemProfiler

Lots of good stuff here, AND its in json!

### hostnames

#### Option 1

nvram -- fmm-computer-name might be there and parseable on NHGRI equipment.
         fmm is "find my mac" and ISNT there on mine.

We can set "nhgri-hostname" and "nhgri-assetno" however...

#### Option 2 (fallback)

Based on <https://itbwiki.nhgri.nih.gov/wiki/index.php/Workstation_Naming_Conventions>

`^HG-(0\d{7})-[DLTVC][WML]\d$`  with gmi flags (global, multiline [^ and $ work], case insensitive)
is the regex for our hostnames.


TODO: make the trailing digit optional - they changed the spec in 2021-12

```python
HG_HOSTNAME_RE = re.compile('^HG-(0\d{7})-[DLTVC][WML](\d)?$', re.ASCII | re.IGNORECASE | re.MULTILINE)
```

## Gentoo

``` text
egglestongc@hg-02224551-vlg ~/Downloads> eselect profile show
Current /etc/portage/make.profile symlink:
  default/linux/amd64/23.0/desktop/systemd
egglestongc@hg-02224551-vlg ~/Downloads> ls -ld /etc/portage/make.profile
lrwxrwxrwx 1 root root 75 Oct 27  2024 /etc/portage/make.profile -> ../../var/db/repos/gentoo/profiles/default/linux/amd64/23.0/desktop/systemd/
egglestongc@hg-02224551-vlg ~/Downloads>

Whatever man, whatever - she says "linux/amd/23.0/desktop/sysstemd"

So... uh?
```

### Reboot detection - Not available in gentoo

### Update detection

`emerge -DupvN @world`

## FreeBSD

### Uname

Its the same

### os.version

freebsd-version
"If no option is specified, it will print the userland version only."

<https://man.freebsd.org/cgi/man.cgi?freebsd-version>
