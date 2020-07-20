package main

import (
	"fmt"
	"github.com/SpencerBrown/get-mongodb/get"
	"github.com/SpencerBrown/get-mongodb/version"
	"log"
	"os/user"
	"path/filepath"
)

const binaryPath = "mongodb-binaries"

func main() {
	myVersion := version.Version{
		Arch:    "x86_64",
		OS:      "macos",
		Distro:  "",
		Release: version.ReleaseType{Version: 4, Major: 2, Minor: 8, Enterprise: true},
	}
	thisUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	myPath := filepath.Join(thisUser.HomeDir, binaryPath)

	myLocation, err := myVersion.ToLocation()
	if err != nil {
		log.Fatalf("Error getting location: %v\n", err)
	}
	myURL := myLocation.URLPrefix + myLocation.Filename
	myPath = filepath.Join(thisUser.HomeDir, binaryPath, myLocation.Filename)
	err = get.DownloadFile(myPath, myURL, 60)
	if err != nil {
		log.Fatalf("Error downloading from URL %s: %v\n", myURL, err)
	} else {
		fmt.Printf("Successfully downloaded %s from URL %s\n", myPath, myURL)
	}

	listVersions()
}

// https://downloads.mongodb.com/osx/mongodb-macos-x86_64-enterprise-4.2.8.tgz

func listVersions() {

}
