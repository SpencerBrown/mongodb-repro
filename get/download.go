package get

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

// Download a file and write it locally
func DownloadFile(filePath string, url string, timeout int) error {

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// Download from URL into memory
	content, err := Download(url, timeout)
	if err != nil {
		return err
	}

	// Write to the file
	contentReader := bytes.NewReader(content)
	_, err = io.Copy(out, contentReader)
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

// Download an archive and expand it to a directory on disk
func DownloadArchive(myPath string, myUrl string, timeout int) error {

	// Download from URL into memory
	content, err := Download(myUrl, timeout)
	if err != nil {
		return err
	}

	// Check type of archive (zip, tgz) and continue accordingly
	parsedURL, err := url.Parse(myUrl)
	if err != nil {
		return err
	}
	fn := path.Base(parsedURL.Path)
	ft := filepath.Ext(fn)
	switch ft {
	case ".zip":
		zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
		if err != nil {
			return err
		}
		for _, f := range zipReader.File {
			thePath := filepath.Join(myPath, f.Name)
			err = os.MkdirAll(filepath.Dir(thePath), 0777)
			uzread, err := f.Open()
			if err != nil {
				return err
			}
			out, err := os.Create(thePath)
			if err != nil {
				return err
			}
			_, err = io.Copy(out, uzread)
			if err != nil {
				return err
			}
			err1 := out.Close()
			err2 := uzread.Close()
			if err1 != nil {
				return err1
			}
			if err2 != nil {
				return err2
			}
			// Set file create/modified time from the zip file
			err = os.Chtimes(thePath, f.Modified, f.Modified)
			if err != nil {
				return err
			}
			// Set file permissions based on the zip file
			err = os.Chmod(thePath, f.Mode())
		}
	case ".tgz":
		gzReader, err := gzip.NewReader(bytes.NewReader(content))
		if err != nil {
			return err
		}
		tarReader := tar.NewReader(gzReader)
		for true {
			tarHeader, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			thePath := filepath.Join(myPath, tarHeader.Name)
			switch tarHeader.Typeflag {
			case tar.TypeDir:
				err = os.Mkdir(thePath, 0777)
				if err != nil {
					return err
				}
			case tar.TypeReg:
				err = os.MkdirAll(filepath.Dir(thePath), 0777)
				out, err := os.Create(thePath)
				if err != nil {
					return err
				}
				_, err = io.Copy(out, tarReader)
				if err != nil {
					return err
				}
				err = out.Close()
				if err != nil {
					return err
				}
				// Set file access/modified time from the tar file
				err = os.Chtimes(thePath, tarHeader.AccessTime, tarHeader.ModTime)
				if err != nil {
					return err
				}
				// Set file permissions based on the zip file
				err = os.Chmod(thePath, tarHeader.FileInfo().Mode())
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("item %s in tar file %s is not a directory or regular file", tarHeader.Name, fn)
			}
		}
		err = gzReader.Close()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("file %s not zip or tgz format", fn)
	}
	return nil
}

// Download file from "url" to memory
func Download(url string, timeout int) ([]byte, error) {

	// Get the data with a custom timeout
	var netClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status '%s' downloading '%s'", resp.Status, url)
	}

	// Copy the data to memory (byte slice)
	content, err := ioutil.ReadAll(resp.Body)
	//_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return content, nil
}
