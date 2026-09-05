# Coding Standards

## Go file structure (subcmd files)
- `Run()` is always the **first** function in a subcmd/* file

## Formatting
- Always add a **blank line after an `if err != nil` block** before the next statement

<!--
## Package layout
- `main.go` at repo root, `package main`
- All other Go files in `internal/`, `package internal`
-->

## CLI (kong)
- Default values belong in kong struct tags (`default:"..."`) not in `Run()`

## Build and Test
- On Darwin, the system `make` is 3.81; use `gmake` instead of `make` for all build and test commands
- Use `make build` to do test builds (preferred over `go build ./...`)
- Use `make test` (via `run_in_terminal`) to run tests — **do not** use the `runTests` tool; it opens VS Code's test UI panel which cannot be closed from chat
- Use `make test-units` when full per-test output is needed for debugging unit tests
- Use `make test-scripts` when full per-test output is needed for debugging testscript tests
- Makefile syntax must remain compatible with **GNU make 3.82** — do not use features introduced after 3.82

## Man pages (*.1.md)
- Use two blank lines before each major `##` section heading
- markdownlint cannot enforce a minimum blank line count — this is a convention only

## Naming
- Variable and function names must be **at least 3 characters** — single-letter and two-letter names are not allowed (except loop indices like `i`, `j` where no semantic meaning is lost)
- Do not give a method and a package-level function the same name — it is valid Go but confuses human readers; ask for clarification if the right name is not obvious

## Communication
- Keep user-facing updates concise: use brief phase preambles, read only the targeted context needed for the immediate change, and summarize results rather than reproducing mechanical edits, raw diffs, or voluminous tool output
- Surface detailed test failures, diagnostics, or logs only when they require a decision or action
- This is an emoji-free project unless the user explicitly requests an emoji to be added or changed

## Editing
- Use editor tools for renames and simple edits — not `sed` or other terminal commands — so the user gets the Zed/VS Code diff/confirmation UX
- Store conventions and preferences in this file (`copilot-instructions.md`), not in Copilot memory — ask before using memory instead

## General
- Always ask for clarification if the intent or requirements are unclear before proceeding
- Code quality priority (in order): **correctness first** (it must be right), **idiomatic Go second** (use language conventions and stdlib patterns), **simplicity third** (prefer the obvious, readable solution over clever or over-engineered ones)

