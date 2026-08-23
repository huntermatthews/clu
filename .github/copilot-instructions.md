# Agent Instructions

## Building

Use `gmake` instead of `make` — the Makefiles require GNU Make 3.82+ which on macOS is provided by `gmake`.

```
gmake build
```

## Sandbox limitation

`go build` (and any `gmake` target that invokes it) will always fail in a sandboxed environment with:

```
go: writing stat cache: open /Users/.../.local/share/go/pkg/mod/cache/... operation not permitted
```

This is a Go bookkeeping write to the module stat cache, not a compilation error. There is no workaround available inside the sandbox. Use the editor's language server diagnostics to verify correctness instead of attempting to build.

Do **not** create temporary directories inside the project worktree as a workaround — that pollutes the git diff.
