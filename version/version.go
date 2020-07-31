package version

import (
	"fmt"
	"regexp"
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
	Modifier   string // e.g. "rc1"
	Enterprise bool   // true if Enterprise edition
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

/*
The following regex breaks down a tarball filename into a slice of strings:
	0. The entire string (filename)
	1. OS
	2. Architecture
	3. "enterprise" or ""
	4. Distro or ""
	5. Release "x.y.z"
	6. Release suffix e.g. "rc14"
	7. ".tgz", ".zip", or ""

	If the filename doesn't match the pattern, an empty slice is returned
*/

const filenameRegexString = `^mongodb-(linux|macos|osx|osx-ssl|win32)-(x86_64|s390x|ppc64le|aarch64)-(?:(enterprise)-)?(?:(rhel\d\d|debian\d\d|suse\d\d|ubuntu\d\d\d\d|amzn64|amazon2|windows-64)-)?(\d\.\d\.\d{1,2})(?:-([a-z0-9]+))?(.tgz|.zip)?`

var filenameRegex *regexp.Regexp

/*
The following regex breaks down a release string into a slice of strings:
	0. The entire string
	1. Release version (single digit string)
	2. Release major (single digit string)
	3. Release minor (one or two-digit string)
	4. Release modifier ("xxx" or "")
*/

const releaseRegexString = `^(\d)\.(\d)\.(\d{1,2})(?:-([a-z0-9]+))?`

var releaseRegex *regexp.Regexp

func init() {
	filenameRegex = regexp.MustCompile(filenameRegexString)
	releaseRegex = regexp.MustCompile(releaseRegexString)
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
	if v.Release.Modifier != "" {
		rel = rel + "-" + v.Release.Modifier
	}
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
	// Filename elements
	/*
			"mongodb-"
			OS-
			Arch-
			"enterprise-" (missing if Community version)
			Distro- (missing for macOS, "windows-64-" for Windows)
			release (e.g. 4.2.8 or 4.4.0-rc14)
			extension (optional) .tgz or .zip

		The regex breaks down a tarball filename into a slice of strings:
		0. The entire string (filename)
		1. OS
		2. Architecture
		3. "enterprise" or ""
		4. Distro or ""
		5. Release "x.y.z"
		6. Release suffix e.g. "rc14"
		7. ".tgz", ".zip", or ""
	*/
	relements := filenameRegex.FindStringSubmatch(fn)
	if len(relements) != 8 {
		return nil, fmt.Errorf("filename '%s' not recognized as a MongoDB tarball", fn)
	}
	thisVersion := new(Version)
	if relements[1] == "osx" || relements[1] == "osx-ssl" {
		thisVersion.OS = "macos"
	} else {
		thisVersion.OS = OSType(relements[1])
	}
	thisVersion.Arch = ArchType(relements[2])
	if relements[3] == "enterprise" {
		thisVersion.Release.Enterprise = true
	}
	thisVersion.Distro = DistroType(relements[4])
	thisVersion.Release.Modifier = relements[6]
	releaseStrings := strings.Split(relements[5], ".")
	thisVersion.Release.Version, _ = strconv.Atoi(releaseStrings[0]) // no error checking necessary,
	thisVersion.Release.Major, _ = strconv.Atoi(releaseStrings[1])   // regex has already vetted the string
	thisVersion.Release.Minor, _ = strconv.Atoi(releaseStrings[2])

	if relements[7] == ".tgz" && thisVersion.OS == "win32" {
		return nil, fmt.Errorf("filename '%s' for Windows requires .zip extension", fn)
	}
	if relements[7] == ".zip" && thisVersion.OS != "win32" {
		return nil, fmt.Errorf("non-Windows filename '%s' requires .tgz extension", fn)
	}

	err := thisVersion.Validate()
	if err != nil {
		return nil, fmt.Errorf("fn '%s' invalid: %v", fn, err)
	}
	return thisVersion, nil
}

// Given a string representing a release, return a ReleaseType
func ToRelease(rs string) (ReleaseType, error) {

	/*
		The following regex breaks down a release string into a slice of strings:
		0. The entire string
		1. Release version (single digit string)
		2. Release major (single digit string)
		3. Release minor (one or two-digit string)
		4. Release modifier ("xxx" or "")
	*/

	rt := ReleaseType{}
	relements := releaseRegex.FindStringSubmatch(rs)
	if len(relements) != 5 {
		return rt, fmt.Errorf("release string '%s' does not match the pattern", rs)
	}
	fmt.Println(relements)
	rt.Version, _ = strconv.Atoi(relements[1]) // don't have to check errors; regex has already vetted the string
	rt.Major, _ = strconv.Atoi(relements[2])
	rt.Minor, _ = strconv.Atoi(relements[3])
	rt.Modifier = relements[4]
	fmt.Println(rt)
	return rt, nil
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
