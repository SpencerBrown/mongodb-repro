package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Arch    ArchType    // hardware platform
	OS      OSType      // operating system
	Distro  DistroType  // Linux distribution
	Release ReleaseType // MongoDB Release
}

type Location struct {
	Filename  string // Filename for directory containing version's download
	URLPrefix string // URL prefix to fetch the tarball from MongoDB
	URLSuffix string // URL suffix (".tgz", ."zip") for the tarball
}

type ReleaseType struct {
	Version    int
	Major      int
	Minor      int
	Enterprise bool // true if Enterprise edition
}

type ArchType string

var validArch = [...]ArchType{"x86_64", "aarch64", "ppc64le", "s390x"}

type OSType string

var validOS = [...]OSType{"linux", "macos", "win32"}

type DistroType string

var validDistro = [...]DistroType{
	// Linux distros
	"amzn64", "amazon2",
	"debian71", "debian81", "debian92", "debian10",
	"rhel57", "rhel62", "rhel67", "rhel70", "rhel71", "rhel72",
	"suse11", "suse12", "suse15",
	"ubuntu1204", "ubuntu1404", "ubuntu1604", "ubuntu1804", "ubuntu2004",
	// macOS distros
	"",
	// Windows distros
	"windows-64",
}

const enterpriseUrlPrefix = "https://downloads.mongodb.com/"
const communityUrlPrefix = "https://fastdl.mongodb.org/"

// Get filename and URL prefix for a Version
func (v *Version) ToLocation() (*Location, error) {
	err := v.Validate()
	if err != nil {
		return nil, err
	}
	rel := fmt.Sprintf("-%d.%d.%d", v.Release.Version, v.Release.Major, v.Release.Minor)
	dist := ""
	if v.Distro != "" {
		dist = "-" + string(v.Distro)
	}
	ent := ""
	if v.Release.Enterprise {
		ent = "-enterprise"
	}

	os := string(v.OS)
	dir := os
	if os == "macos" {
		dir = "osx"
		if (v.Release.Version < 4) || (v.Release.Major < 2) {
			os = "osx-ssl"
			if v.Release.Enterprise {
				os = "osx"
			}
		}
	}
	fn := "mongodb-" + os + "-" + string(v.Arch) + ent + dist + rel

	prefix := communityUrlPrefix
	if v.Release.Enterprise {
		prefix = enterpriseUrlPrefix
	}
	// Suffix is .tgz except for Windows which is .zip
	suffix := ".tgz"
	if v.OS == "win32" {
		suffix = ".zip"
	}
	return &Location{
		Filename:  fn,
		URLPrefix: prefix + dir + "/",
		URLSuffix: suffix,
	}, nil
}

// Convert filename to Version
func ToVersion(fn string) (*Version, error) {
	// Filename elements ,separated by dashes:
	/*
		"mongodb"
		OS
		Arch
		"enterprise" (missing if Community version)
		Distro (missing for macOS, "windows-64" for Windows)
		release (e.g. 4.2.8.tgz; will also work without the extension)
	*/
	elements := strings.Split(fn, "-")
	n := len(elements)
	if (n < 4) || (n > 7) {
		return nil, fmt.Errorf("fn '%s' invalid, too few or too many elements", fn)
	}
	if elements[0] != "mongodb" {
		return nil, fmt.Errorf("fn '%s' invalid, does not start with 'mongodb'", fn)
	}
	elements = elements[1:]
	thisVersion := new(Version)
	// before 4.2, macOS was "osx" in the name for Enterprise, and "osx-ssl" for Community
	if elements[0] == "osx" {
		thisVersion.OS = "macos"
		if elements[1] == "ssl" {
			elements = elements[1:]
		}
	} else {
		thisVersion.OS = OSType(elements[0])
	}
	thisVersion.Arch = ArchType(elements[1])
	elements = elements[2:]
	if elements[0] == "enterprise" {
		thisVersion.Release.Enterprise = true
		if len(elements) <= 1 {
			return nil, fmt.Errorf("fn '%s' invalid, nothing after 'enterprise'", fn)
		}
		elements = elements[1:]
	}
	if thisVersion.OS == "linux" {
		thisVersion.Distro = DistroType(elements[0])
		if len(elements) != 2 {
			return nil, fmt.Errorf("fn '%s' invalid, nothing after Linux distro", fn)
		}
		elements = elements[1:]
	}
	if thisVersion.OS == "win32" {
		// Distro must be "windows", "64"
		if (len(elements) != 3) || (elements[0] != "windows") || (elements[1] != "64") {
			return nil, fmt.Errorf("fn '%s' invalid, Windows must be windows-64 distro", fn)
		}
		thisVersion.Distro = "windows-64"
		elements = elements[2:]
	}
	rel := strings.Split(elements[0], ".") // "x.y.z.ft" or "x.y.z"
	if len(rel) == 4 {
		ft := rel[3]
		rel = rel[:3]
		if thisVersion.OS == "win32" {
			if ft != "zip" {
				return nil, fmt.Errorf("fn '%s': Windows filetype must be 'zip'", fn)
			}
		} else {
			if ft != "tgz" {
				return nil, fmt.Errorf("fn '%s': non-Windows filetype must be 'tgz'", fn)
			}
		}
	}
	if len(rel) != 3 {
		return nil, fmt.Errorf("fn '%s' invalid, release must have 3 elements", fn)
	}
	var r [3]int
	var err error
	for i := 0; i < 3; i++ {
		r[i], err = strconv.Atoi(rel[i])
		if err != nil {
			return nil, fmt.Errorf("fn '%s' invalid, release must be numeric", fn)
		}
	}
	thisVersion.Release.Version = r[0]
	thisVersion.Release.Major = r[1]
	thisVersion.Release.Minor = r[2]
	err = thisVersion.Validate()
	if err != nil {
		return nil, fmt.Errorf("fn '%s' invalid: %v", fn, err)
	}
	return thisVersion, nil
}

// Validate a Version
func (v *Version) Validate() error {
	invalid := true
	for _, x := range validArch {
		if x == v.Arch {
			invalid = false
		}
	}
	if invalid {
		return fmt.Errorf("%s is not a valid architecture", v.Arch)
	}
	invalid = true
	for _, x := range validOS {
		if x == v.OS {
			invalid = false
		}
	}
	if invalid {
		return fmt.Errorf("%s is not a valid operating system", v.OS)
	}
	invalid = true
	for _, x := range validDistro {
		if x == v.Distro {
			invalid = false
		}
	}
	if invalid {
		return fmt.Errorf("%s is not a valid distribution", v.Distro)
	}
	if (v.Release.Version < 2) || (v.Release.Version > 4) {
		return fmt.Errorf("release Version %d must be 2, 3, or 4", v.Release.Version)
	}
	if (v.Release.Major < 0) || (v.Release.Major > 6) {
		return fmt.Errorf("major Release %d must be 0 through 6", v.Release.Major)
	}
	if (v.Release.Minor < 0) || (v.Release.Minor > 40) {
		return fmt.Errorf("minor Release %d must be 0 through 40", v.Release.Minor)
	}
	return nil
}
