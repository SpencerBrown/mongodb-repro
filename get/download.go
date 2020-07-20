package get

import (
	"io"
	"net/http"
	"os"
	"time"
)

// Download a file and write it locally
func DownloadFile(filepath string, url string, timeout int) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// Download the file
	err = Download(out, url, timeout)
	if err != nil {
		return err
	}

	// Close the file
	err = out.Close()
	if err != nil {
		return err
	}

	return nil
}

// Download file from "url" to "out" with timeout
func Download(out io.Writer, url string, timeout int) error {

	// Get the data with a custom timeout
	var netClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return err
	}

	// Write the body to the writer
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	return err
}
