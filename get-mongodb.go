package main

import (
	"fmt"
	"github.com/SpencerBrown/get-mongodb/version"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	myVersion := version.Version{
		Arch:    "x86_64",
		OS:      "macos",
		Distro:  "",
		Release: version.ReleaseType{Version: 4, Major: 2, Minor: 8, Enterprise: true},
	}
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	myLocation := filepath.Join(user.HomeDir, "mongodb-binaries")

	URL, err := myVersion.ToURL()
	if err != nil {
		log.Fatalf("Error getting URL: %v\n", err)
	}
	fn, err := myVersion.ToFilename()
	if err != nil {
		log.Fatalf("Error getting filename: %v\n", err)
	}
	err = DownloadFile(filepath.Join(myLocation, fn), URL)
	if err != nil {
		log.Fatalf("Error downloading from URL %s: %v\n", URL, err)
	} else {
		fmt.Printf("Successfully downloaded from URL %s\n", URL)
	}
}

// https://downloads.mongodb.com/osx/mongodb-macos-x86_64-enterprise-4.2.8.tgz

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
