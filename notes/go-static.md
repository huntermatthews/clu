# Go Static Builds

Go doesn't always build static - usage of net! or os/user and other things can trigger a dl-link against libc on linux.

Various hacks include:

- `CGO_ENABLED=0 go build ...`
- Build tags
  - `go build -tags osusergo`
  - `go build -tags netgo`
  - `go build -tags osusergo,netgo`
- Linker flags
  - go build -ldflags="-extldflags=-static"
  - This works (maybe) when you want Sqlite or something inherently C
- More linker magic:
  - `go build -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build'`
