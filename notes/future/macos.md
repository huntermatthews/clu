# MacOS

## Uptime

Parsing uptime(4) has a regex that is unreliable and goofy - plus parsing the result.
`sysctl kern.boottime` is far more consistent.
