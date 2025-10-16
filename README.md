# Clu

Clu gives us a clue as to what this os / hardware is.

## Development Notes

This program grew out of multiple generations of shell scripting (and one previous disastrous attempt at python) -- keep that in mind
before judging too harshly.

### Layout Notes

- Whole program (script) was originally just the "report" subcmd.
- clu/cmd should be just the command line handling code  - but older subcmds like report and archive are all in one file. Remember - shell.
- functions named do_* should be the main target of a command line subcmd.
- Some of the parsers in sources/* could use some love.
