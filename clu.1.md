<!-- markdownlint-disable -->
<!--
Copyright (c) 2011 Hunter Matthews <hunter@unix.haus>
SPDX-License-Identifier: GPL-2.0-only

The linter is disabled because thats for generic markdown files and this needs to be man page / go-md2man specific.
-->

clu 1 "January 2011" clu "NHGRI Sysadmin Commands"
============================================

## NAME

clu - Describes various interesting facts about the current system.

## SYNOPSIS

**clu** [**--help** | **--version** | **--debug** | **--net**]

**clu** **collector**

**clu** **facts** [**fact-names** ...] [**--tier**, **-t** {**1** | **2** | **3**}] [**--out** {**dots** | **shell** | **json**}]

**clu** **requires** {**list** | **check**}

## DESCRIPTION

**clu** is a tool for reporting various facts about a system. It reads well known system files and programs to collect the information. It was created because other tools in this category (and there are *many*) either tried to report too much (overwhelming a busy sysadmin) or far too little (causing the sysadmin to have to run other commands for a complete picture). It is therefore strongly opinionated - what it outputs and how that is formatted are pretty much at the opinion of the author.

## OPTIONS

The following global options are supported. Individual commands may have their own specific options.

#### --help, -h

List the help screen for the command and exit.

#### --version

Print the version of the command and exit.

#### --debug

Print extra debugging output (not really implemented yet).

#### --net

Allow access to the network. Some steps require more than just a DNS query (such as talking to Device42 inventory system or updating the DNF/APT data - this flag allows such access).

## COMMANDS

#### collector

Create a compressed **tar**(1) archive of all the required file and program outputs used by **clu**. This facilitates both testing and debugging without needing to be fully networked all the time or remote debugging. The archive will be named /tmp/clu_\<hostname\>.tgz.

#### facts [fact-names ...] [--tier, -t {1 | 2 | 3}] [--out {dots | shell | json}]

Report system facts. Without arguments, reports all facts for the specified tier.

**fact-names**: Optional space-separated list of specific facts to report. Use `clu requires list` to see available fact names.

**--tier** (*default: 1*): Detail level where higher numbers provide more information.

**--out** (*default: dots*): Output format.

#### requires {check | list}

List or check the required files, programs and api's required for the program to run.

## EXAMPLES

Check what system facts are available:
```
$ clu requires list
```

Verify all system requirements are met:
```
# clu requires check
```

Generate a quick system overview:
```
$ clu facts --tier 1
```

Create a detailed system report for troubleshooting:
```
$ clu facts --tier 3 --out json > system-report.json
```

Collect diagnostic data for remote analysis:
```
# clu collector
# scp /tmp/clu_$(hostname).tgz support@example.com:
```

Export system facts as shell variables:
```
$ eval "$(clu facts --out shell)"
$ echo "CPU: $CPU_MODEL, OS: $OS_NAME"
```

## EXIT STATUS

**clu** exits with one of the following values:

    0    Success.

    1    General error or invalid usage.

    2    System requirements not met (for `requires check`).

## NOTES

**clu** attempts to gracefully handle missing system files or programs. Use `clu requires check` to verify system compatibility.

Network access (via `--net`) is only required for certain advanced facts that query external services.

Some facts require root privileges to access system files in `/proc`, `/sys`, or hardware interfaces.


## SEE ALSO

**lscpu**(1), **lshw**(1), **dmidecode**(8), **hostnamectl**(1), **systemd-detect-virt**(1), **uname**(1), **free**(1), **df**(1), **ip**(8), **ss**(8), **tar**(1)


## HISTORY

**clu** is the 4th or 5th incarnation of this idea, going back to 2011. I finally gave up on 2 different incarnations of shell script. The deployment is easier but the programming is AWFUL. Pure python zipapps solve the deployment problem (mostly). (I wish Python had the same deployment model of Go...) Names have varied considerably over the years. Earlier incarnations strongly influenced the `collector` command and being able to test the parsers without being logged in to networked systems.

## AUTHORS

Hunter Matthews (hunter@unix.haus) is to blame on this one.

<!--
ALTERNATIVE COMMANDS SECTION FROM V2 (more detailed for complex subcommands):

## COMMANDS

### collector

Create a compressed **tar**(1) archive of all the required file and program outputs used by **clu**. This facilitates both testing and debugging without needing to be fully networked all the time or remote debugging. The archive will be named `/tmp/clu_<hostname>.tgz`.

### facts [OPTIONS]

Report various facts about the computer. This includes details of the OS, hardware (virtual or physical), networking, etc.

**Command-specific options:**

**--tier**, **-t**=*TIER*
: Controls tier of facts (more and more details at higher / larger number tiers...) are printed.

: Valid values: **1**, **2**, **3**

: The default is **1**.

**--out**=*FORMAT*
: Output format for the facts.

: Valid values: **dots**, **shell**, **json**

: The default is **dots**.

### requires {check | list}

List or check the required files, programs and api's required for the program to run.
-->
