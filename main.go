package main

import (
	"flag"
	"fmt"
	"github.com/SpencerBrown/mongodb-repro/cmds"
	"github.com/SpencerBrown/mongodb-repro/version"
	"os"
)

type OptionsType struct {
	Arch      *string
	MyOS      *string
	Distro    *string
	Release   *string
	Community *bool
}

func printHelp() {
	fmt.Printf("%s list - lists currently downloaded versions\n", os.Args[0])
	fmt.Printf("%s get - downloads a version\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	opts := OptionsType{
		Arch:      flag.String("arch", "x86_64", "Architecture: x86_64 (default), aarch64, ppc64le, s390x"),
		MyOS:      flag.String("os", "linux", "OS: linux (default), macos, win32"),
		Distro:    flag.String("distro", "ubuntu1604", "Linux distro"),
		Release:   flag.String("version", "4.2.9", "MongoDB version"),
		Community: flag.Bool("community", false, "Community version?"),
	}
	flag.Parse()

	if flag.NArg() != 1 {
		printHelp()
		return
	}

	v := &version.Version{
		Arch:   version.ArchType(*opts.Arch),
		OS:     version.OSType(*opts.MyOS),
		Distro: version.DistroType(*opts.Distro),
	}
	var err error
	v.Release, err = version.ToRelease(*opts.Release)
	if err != nil {
		fmt.Printf("Error in release '%s': %v\n", *opts.Release, err)
		return
	}
	v.Release.Enterprise = !(*opts.Community)
	err = v.Validate()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var isWindows = v.OS == "win32"

	err = cmds.Cmds(flag.Arg(0), v, isWindows)
	if err != nil {
		fmt.Printf("Error processing request: %v", err)
	}
}
