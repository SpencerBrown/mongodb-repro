package main

import (
	"fmt"
	"github.com/SpencerBrown/get-mongodb/get"
	"github.com/SpencerBrown/get-mongodb/version"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
)

const binaryDir = "mongodb-binaries"

func getPath() string {
	thisUser, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %v\n", err)
	}
	return filepath.Join(thisUser.HomeDir, binaryDir)
}

var binaryPath = getPath()

func main() {
	err := getOneAndExpand()
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	err = listVersions()
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
}

func getOneAndExpand() error {
	myVersion := version.Version{
		Arch:    "x86_64",
		OS:      "macos",
		Distro:  "",
		Release: version.ReleaseType{Version: 4, Major: 0, Minor: 18, Enterprise: false},
	}
	myLocation, err := myVersion.ToLocation()
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
	myVersion2, err := version.ToVersion(myLocation.Filename)
	fmt.Println(myVersion2)

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
		fmt.Printf("%-7s %-5s %-10s %d.%d.%2d %s\n", v.Arch, v.OS, v.Distro, v.Release.Version, v.Release.Major, v.Release.Minor, isEnterprise)
	}
	return nil
}
