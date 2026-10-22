# SBOM generation options

```shell
SYFT_FORMAT_SPDX_JSON_AUTHOR="Organization: YourCompany LLC (security@yourcompany.com)" \
SYFT_FORMAT_SPDX_JSON_NAMESPACE="https://yourcompany.com" \
syft file:./my-go-binary \
  -o spdx-json=sbom.spdx.json
```
