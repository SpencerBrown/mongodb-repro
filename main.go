package main

import (
	"flag"
	"fmt"
	"github.com/SpencerBrown/content"
	"github.com/SpencerBrown/mongodb-repro/cmds"
	"github.com/SpencerBrown/mongodb-repro/version"
	"github.com/SpencerBrown/mongodb-repro/web"
	"log"
	"os"
	"path/filepath"
)

type OptionsType struct {
	Arch      *string
	OS        *string
	Distro    *string
	Release   *string
	Community *bool
	UI        *bool
}

func printHelp() {
	fmt.Printf("%s list - lists currently downloaded versions\n", os.Args[0])
	fmt.Printf("%s get - downloads a version\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	// Generate static content if "go generate" was run
	if len(os.Args) > 1 && os.Args[1] == "generate" {
		// Handle generating static assets
		log.Println("Generating static assets as Go source code...")
		err := generateStatic()
		if err != nil {
			log.Fatalf("Error generating static assets: %v\n", err)
		}
		log.Println("Done generating Go source code for static assets")
		return
	}

	opts := OptionsType{
		Arch:      flag.String("arch", "x86_64", "Architecture: x86_64 (default), aarch64, ppc64le, s390x"),
		OS:        flag.String("os", "linux", "OS: linux (default), macos, win32"),
		Distro:    flag.String("distro", "ubuntu1604", "Linux distro"),
		Release:   flag.String("version", "4.2.9", "MongoDB version"),
		Community: flag.Bool("community", false, "Community version?"),
		UI:        flag.Bool("ui", false, "Invoke Web UI?"),
	}
	flag.Parse()

	if *opts.UI {
		web.Server("localhost", 8888, true)
		return
	}

	if flag.NArg() != 1 {
		printHelp()
		return
	}

	v := &version.Version{
		Arch:   version.ArchType(*opts.Arch),
		OS:     version.OSType(*opts.OS),
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

// create the static.go file from the static content files in the "static-content" directory via "go:generate"
// put static.go in the "staticContent" directory and give it the package name "static-content"
//go:generate mongodb-repro generate staticContent
func generateStatic() error {

	const staticOutputName = "static.go"
	const staticContentDir = "staticContent"

	staticFilesDir, err := filepath.Abs(staticContentDir)
	if err != nil {
		return fmt.Errorf("unable to locate static files: %v", err)
	}

	err = content.GenerateContent(staticFilesDir, staticFilesDir, staticOutputName, staticContentDir)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}
