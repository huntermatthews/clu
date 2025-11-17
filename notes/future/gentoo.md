
# Gentoo

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

## Reboot detection - Not available in gentoo

## Update detection

`emerge -DupvN @world`
