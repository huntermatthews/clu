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

## Security & Verification

All release artifacts include SLSA Level 3 attestations for supply chain verification. You can verify the authenticity and integrity of any downloaded artifact:

### Verify Raw Binaries

```bash
gh attestation verify clu-linux-amd64 --owner NHGRI
gh attestation verify clu-linux-arm64 --owner NHGRI
gh attestation verify clu-darwin-arm64 --owner NHGRI
```

### Verify Package Files

```bash
# RPM packages
gh attestation verify clu-1.0.0-1.x86_64.rpm --owner NHGRI
gh attestation verify clu-1.0.0-1.aarch64.rpm --owner NHGRI
gh attestation verify clu-1.0.0-1.x86_64.el7.rpm --owner NHGRI

# DEB packages
gh attestation verify clu_1.0.0-1_amd64.deb --owner NHGRI
gh attestation verify clu_1.0.0-1_arm64.deb --owner NHGRI
```

### Verify Checksums

```bash
# Download and verify checksums
curl -fsSL -O https://github.com/NHGRI/clu/releases/download/v1.0.0/SHA256SUMS
curl -fsSL -O https://github.com/NHGRI/clu/releases/download/v1.0.0/clu-linux-amd64
sha256sum -c SHA256SUMS
```


## Private repo download information

```shell
DOWNLOAD_URL="PASTE_THE_browser_download_url_HERE"
FILE_NAME="name_of_your_asset.tar.gz"

curl -L \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Accept: application/octet-stream" \
  -o "${FILE_NAME}" \
  "${DOWNLOAD_URL}"
```
