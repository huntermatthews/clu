# Testing locally built rpms

I’m using itbshell01 for the main build
HOWEVER, thats too new for our EL7 systems. (different versions of rpm)


## 1. Copy rpm to itbsalt03:/tmp

## 2. ssh to itbsalt03, become root

`ssh itbsalt03`
`sudo -i`

## 3. Copy rpm to all the targets for testing

This is about a 1.5MB rpm so this is doable with salt-cp. Much over 4MB really isn’t. Note the -C is required for non-text files.

`salt-cp -L nhgrioradev03,itbsaltdev03,itbsaltdev08,biowebdev05,itbsaltdev06,itbjenkinsdev01,itbredhat01,itbsaltdev07 -C /tmp/clu-v1.0.0+foo-1.el9.x86_64.rpm /tmp/clu-v1.0.0+foo-1.el9.x86_64.rpm`

## 4. Attempt install of the rpm

`salt -L nhgrioradev03,itbsaltdev03,itbsaltdev08,biowebdev05,itbsaltdev06,itbjenkinsdev01,itbredhat01,itbsaltdev07 cmd.run 'rpm -Uvh /tmp/clu-v1.0.0+foo-1.el9.x86_64.rpm’`

```shell
biowebdev05:
    error: Failed dependencies:
     rpmlib(PayloadIsZstd) <= 5.4.18-1 is needed by clu-v1.0.0+foo-1.el9.x86_64
```

All the others succeeded.

## 5. Run the program

```salt -L nhgrioradev03,itbsaltdev03,itbsaltdev08,biowebdev05,itbsaltdev06,itbjenkinsdev01,itbredhat01,itbsaltdev07 cmd.run 'which clu; clu --version ; clu -t3'```

The version is unknown but otherwise it looks ok (for now)

## 6. Cleanup - remove installed rpm and remove the file from /tmp

```salt -L nhgrioradev03,itbsaltdev03,itbsaltdev08,biowebdev05,itbsaltdev06,itbjenkinsdev01,itbredhat01,itbsaltdev07 cmd.run 'rpm -e clu ; rm -f /tmp/clu-v*’```

## Biowebdev05

```shell
[matthewsht@biowebdev05]~/projects% git --version
git version 1.8.3.1

[matthewsht@biowebdev05]~/projects/clu% scripts/build.sh
go: downloading github.com/alecthomas/kong v1.13.0
github.com/NHGRI/clu/pkg: cannot compile Go 1.25 code
github.com/NHGRI/clu/pkg/facts/types: cannot compile Go 1.25 code

Changed the version in go.mod
Compiled clean.

Installing rpmbuild and rpmdevtools. (rpm-build and rpmdev-setuptree respectively)
dnf install -y rpm-build rpmdevtools
```
