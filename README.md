# Clu

Clu gives us a clue as to what this os / hardware is. Please see the man page for more details.

Development targets are in the Justfile.

## Package Versioning Compatibility

When creating releases for both RPM and DEB packages, use version strings that are compatible with both package managers:

### ✅ Safe/Compatible Version Patterns

- **Simple semantic versions:** `1.2.3`, `2.0.0`, `0.1.5`
- **With revision:** `1.2.3-1` (becomes release in RPM, debian_revision in DEB)
- **Pre-release versions:** `1.2.3-alpha`, `1.2.3-beta`, `1.2.3-rc1`, `1.2.3-dev`, `1.2.3-pre1`
- **Date-based versions:** `1.2.3-20241230`, `2024.12.30`

### ❌ Avoid These Patterns

- **Tilde versions:** `1.2.3~alpha` (different sort behavior)
- **Epoch/Colon versions:** `1:1.2.3` (illegal in RPM)
- **Plus versions:** `1.2.3+build1` (inconsistent behavior)

**Current approach:** Git tags like `v1.2.3` are stripped to `1.2.3` for packaging, which is optimal for both RPM and DEB compatibility.
