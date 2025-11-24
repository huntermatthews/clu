package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	version string = "0.10"
)

var (
	mockFlag    string
	versionFlag bool
)


func cmdlineParse() {
	flag.StringVar(&mockFlag, "mock", "host1", "mock help message")
	flag.BoolVar(&versionFlag, "version", false, "version help message")

	flag.Parse()

	// TODO: hide behind --debug flag
	fmt.Println("debug: mockFlag value is:", mockFlag)
	fmt.Println("debug: versionFlag value is:", versionFlag)

}

func main() {

	cmdlineParse()

	if versionFlag {
		fmt.Println("clu v", version)
		os.Exit(0)
	}

	os.Exit(0)
}
