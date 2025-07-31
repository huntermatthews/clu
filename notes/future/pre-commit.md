# pre-commit stuff

```yaml
fail_fast: true
repos:
  - repo: local
    hooks:
      - id: ruff-check
        name: Ruff check
        entry: poetry run ruff check
        args: [--fix]
        language: system
        types: [file, python]
      - id: ruff-format
        name: Ruff format
        entry: poetry run ruff format
        language: system
        types: [file, python]
      - id: pyright
        name: Pyright type check
        entry: poetry run pyright
        language: system
        types: [file, python]
```

```yaml
repos:
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: "v0.7.4"
    hooks:
      - id: ruff
        args: ["--fix", --max-line-length=120]
    - id: ruff-format
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v5.0.0
    hooks:
      - id: check-yaml
```
