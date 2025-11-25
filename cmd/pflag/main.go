package main

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
)

var ip = flag.Int("flagname", 1234, "help message for flagname")

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
