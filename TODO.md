# TODO

## Phase One scope decision

- Do not move the shared `Clu` fact source into a `pkg/facts/generic/` package. This is no longer a project goal.

## Remaining work

- [ ] Complete `LinuxCheckUpdate.Requires()` so `collector` captures the package-manager commands used by the update-check source.
- [ ] Remap the `TestRun_Hosts` output fixtures. This is in progress and is not a Phase One goal.
- [ ] Add the missing `ipmitool` mock fixture for `testdata/host1` so collector input fixtures are complete.
- [ ] Defer collector archive-content assertions for now.
- [ ] Ignore the `idea-register3.md` inactive-checkupdate filename cleanup for now.
