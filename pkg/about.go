package pkg

// Version of the program, set at build time -- consts won't work with -ldflags
var Version = "unknown"

// About the program itself
const (
	Title       = "clu"
	Summary     = "clu utility"
	License     = "Apache-2.0"
	AuthorEmail = "hunter@unix.haus"
	AuthorName  = "Hunter Matthews"
)

// func reportBuildInfo() {
// 	info, ok := debug.ReadBuildInfo()
// 	if !ok {
// 		fmt.Println("Build info not available")
// 		return
// 	}

// 	fmt.Printf("Build Information:\n")
// 	fmt.Printf("  GoVersion: %s\n", info.GoVersion)
// 	fmt.Printf("  Path: %s\n", info.Path)
// 	fmt.Printf("  Main Module: %s@%s\n", info.Main.Path, info.Main.Version)
// 	for count, dep := range info.Deps {
// 		fmt.Printf("  Dependency[%d]: %s@%s\n", count, dep.Path, dep.Version)
// 	}
// 	for _, setting := range info.Settings {
// 		fmt.Printf("  %s: %s\n", setting.Key, setting.Value)
// 	}
// }
