package global

import (
	"fmt"
	"runtime/debug"
)

// Version of the program, set at build time -- consts won't work with -ldflags
var Version = "unknown"

// About the program itself
const (
	Title       = "clu"
	Summary     = "clu utility"
	License     = "Non-Distributable"
	AuthorEmail = "hunter@unix.haus"
	AuthorName  = "Hunter Matthews"
)

func reportBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Build info not available")
		return
	}

	fmt.Printf("Build Information:\n")
	fmt.Printf("  GoVersion: %s\n", info.GoVersion)                           // this
	fmt.Printf("  Path: %s\n", info.Path)                                     // not this field
	fmt.Printf("  Main Module: %s @ %s\n", info.Main.Path, info.Main.Version) // this

	fmt.Printf("\n\nDependencies:\n")
	for count, dep := range info.Deps {
		fmt.Printf("  Dependency[%d]: %s@%s\n", count, dep.Path, dep.Version)
	}

	fmt.Printf("\n\nSettings:\n")
	for _, setting := range info.Settings {
		fmt.Printf("  %s: %s\n", setting.Key, setting.Value)
	}
}

// goversion   string
// main module  string
// main version string. **** parse this!!
// deps []string of dep.Path + "@" + dep.Version.   == (devel)
// settings map[string]string of setting.Key -> setting.Value
// CGO_ENABLED:
// GOOS:
// GOARCH:
// vcs: git.   optional
// vcs.revision: 67399a5b661390d23e1c556ae8e8c7e89146a0eb opt
// vcs.time: 2026-01-08T04:11:41Z.  opt
// vcs.modified: true/false.  opt

// Build Information:
//   GoVersion: go1.25.3
//   Path: github.com/huntermatthews/clu/cmd2
//   Main Module: github.com/huntermatthews/clu @ v0.8.30-0.20260108041141-67399a5b6613+dirty
// Dependencies:
//
// Settings:
//   -buildmode: exe
//   -compiler: gc
//   -ldflags: -X main.version=99.99.99
//   DefaultGODEBUG: asynctimerchan=1,containermaxprocs=0,decoratemappings=0,gotestjsonbuildtext=1,gotypesalias=0,httpcookiemaxnum=0,httplaxcontentlength=1,httpmuxgo121=1,httpservecontentkeepheaders=1,multipathtcp=0,panicnil=1,randseednop=0,rsa1024min=0,tls10server=1,tls3des=1,tlsmlkem=0,tlsrsakex=1,tlssha1=1,tlsunsafeekm=1,updatemaxprocs=0,winreadlinkvolume=0,winsymlink=0,x509keypairleaf=0,x509negativeserial=1,x509rsacrt=0,x509sha256skid=0,x509usepolicies=0
//   CGO_ENABLED: 1
//   CGO_CFLAGS:
//   CGO_CPPFLAGS:
//   CGO_CXXFLAGS:
//   CGO_LDFLAGS:
//   GOARCH: arm64
//   GOOS: darwin
//   GOARM64: v8.0
//   vcs: git
//   vcs.revision: 67399a5b661390d23e1c556ae8e8c7e89146a0eb
//   vcs.time: 2026-01-08T04:11:41Z
//   vcs.modified: true

// -- about as a struct - like a sane person.
// var (
// 	Version   string
// 	BuildTime string
// 	GitCommit string
// )

// The Config struct
// type BuildInfo struct {
// 	Version   string
// 	BuildTime string
// 	GitCommit string
// }

// func GetBuildInfo() BuildInfo {
// 	return BuildInfo{
// 		Version:   Version,
// 		BuildTime: BuildTime,
// 		GitCommit: GitCommit,
// 	}
// }

// Example of an anonymous struct
// person := struct {
//         Name string
//         Age  int
//     }{
//         Name: "Alice",
//         Age:  30,
//     }
