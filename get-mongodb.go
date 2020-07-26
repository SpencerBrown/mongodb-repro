package main

import (
	"flag"
	"fmt"
	"github.com/SpencerBrown/get-mongodb/get"
	"github.com/SpencerBrown/get-mongodb/version"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const binaryDir = "mongodb-binaries"

func getPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v\n", err)
	}
	return filepath.Join(homedir, binaryDir)
}

var binaryPath = getPath()

func main() {

	var arch = flag.String("arch", "x86_64", "Architecture: x86_64 (default), aarch64, ppc64ld, s390x")
	//var a_arch = flag.String("a", "x86_64", "(short for 'arch')")
	var myos = flag.String("os", "linux", "OS: linux (default), macos, win32")
	//var a_myos = flag.String("o", "linux", "(short for 'os')")
	var distro = flag.String("distro", "ubuntu1604", "Linux distro")
	var release = flag.String("version", "4.2.9", "MongoDB version")
	var community = flag.Bool("community", false, "Community version?")

	flag.Parse()

	if flag.NArg() != 1 {
		printHelp()
		return
	}
	switch flag.Arg(0) {
	case "list":
		err := listVersions()
		if err != nil {
			fmt.Printf("Error listing versions: %v\n", err)
		}
	case "get":
		v := &version.Version{
			Arch:   version.ArchType(*arch),
			OS:     version.OSType(*myos),
			Distro: version.DistroType(*distro),
		}
		var err error
		v.Release, err = version.ToRelease(*release)
		if err != nil {
			fmt.Printf("Error in release '%s': %v\n", *release, err)
			return
		}
		v.Release.Enterprise = !(*community)
		fmt.Println(v)
		err = v.Validate()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		err = getOneAndExpand(v)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Printf("%s list - lists currently downloaded versions\n", os.Args[0])
	fmt.Printf("%s get - downloads a version\n", os.Args[0])
	flag.PrintDefaults()

}

//func main() {
//	err := getOneAndExpand()
//	if err != nil {
//		fmt.Printf("error %v", err)
//		return
//	}
//	err = listVersions()
//	if err != nil {
//		fmt.Printf("error %v", err)
//		return
//	}
//}

func getOneAndExpand(v *version.Version) error {
	myLocation, err := v.ToLocation()
	if err != nil {
		return fmt.Errorf("Error getting location: %v\n", err)
	}
	myURL := myLocation.URLPrefix + myLocation.Filename + myLocation.URLSuffix
	//myPath = filepath.Join(thisUser.HomeDir, binaryPath, myLocation.Filename)
	err = get.DownloadArchive(binaryPath, myURL, 60)
	if err != nil {
		return fmt.Errorf("Error downloading from URL %s: %v\n", myURL, err)
	} else {
		fmt.Printf("Successfully downloaded to %s from URL %s\n", binaryPath, myURL)
	}
	//myVersion2, err := version.ToVersion(myLocation.Filename)
	//fmt.Println(myVersion2)
	return nil
}

// https://downloads.mongodb.com/osx/mongodb-macos-x86_64-enterprise-4.2.8.tgz

func listVersions() error {

	// List available files/directories
	files, err := ioutil.ReadDir(binaryPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(files); i++ {
		fn := files[i].Name()
		v, err := version.ToVersion(fn)
		if err != nil {
			return err
		}
		isEnterprise := "Community"
		if v.Release.Enterprise {
			isEnterprise = "Enterprise"
		}
		fmt.Printf("%-7s %-5s %-10s %d.%d.%2d %s %s\n", v.Arch, v.OS, v.Distro, v.Release.Version, v.Release.Major, v.Release.Minor, v.Release.Modifier, isEnterprise)
	}
	return nil
}
