package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	version string = "0.10"
)

var (
	mockFlag    string
	versionFlag bool
)

type Osrel struct {
	Version string
	Distro  string
}

func cmdlineParse() {
	flag.StringVar(&mockFlag, "mock", "host1", "mock help message")
	flag.BoolVar(&versionFlag, "version", false, "version help message")

	flag.Parse()

	// TODO: hide behind --debug flag
	fmt.Println("debug: mockFlag value is:", mockFlag)
	fmt.Println("debug: versionFlag value is:", versionFlag)

}

func readFile(pwd string, mockFlag string) Osrel {

	file, err := os.Open(filepath.Join(pwd, "tests/mock_data", mockFlag, "/etc/os-release"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res Osrel
	for scanner.Scan() {
		// fmt.Println("debug: ", scanner.Text())
		key, value := parseLine(scanner.Text())
		switch key {
		case "ID":
			res.Distro = value
		case "VERSION_ID":
			res.Version = value
		}

	}
	return res

	// if err := scanner.Err(); err != nil {
	// 	panic(err)
	// }
}

func parseLine(raw string) (key string, value string) {
	//text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	text := raw
	keyValue := strings.Split(text, "=")
	// BUG: make sure we got two and ONLY two parts.
	return keyValue[0], keyValue[1]

}

func main() {

	cmdlineParse()

	if versionFlag {
		fmt.Println("clu v", version)
		os.Exit(0)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("pwd is", pwd)

	res := readFile(pwd, mockFlag)
	fmt.Println("distro:", res.Distro)
	fmt.Println("version:", res.Version)

	os.Exit(0)
}
