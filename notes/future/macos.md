# MacOS

## Uptime

Parsing uptime(4) has a regex that is unreliable and goofy - plus parsing the result.
`sysctl kern.boottime` is far more consistent.

## Ram

❯ sysctl hw.memsize
hw.memsize: 34359738368

## SystemProfiler

Lots of good stuff here, AND its in json!
