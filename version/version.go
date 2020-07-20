package version

import "fmt"

type Version struct {
	Arch    ArchType    // hardware platform
	OS      OSType      // operating system
	Distro  DistroType  // Linux distribution
	Release ReleaseType // MongoDB Release
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

var validOs = [...]OSType{"linux", "macos", "win32"}

type DistroType string

var validDistro = [...]DistroType{
	// Linux distros
	"anzn64", "amazon2",
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

// Get filename for a Version
func (v *Version) ToFilename() (string, error) {
	err := v.Validate()
	if err != nil {
		return "", err
	}
	// Filetype is tgz except for Windows which is zip
	ft := "tgz"
	if v.OS == "win32" {
		ft = "zip"
	}
	rel := fmt.Sprintf("-%d.%d.%d.%s", v.Release.Version, v.Release.Major, v.Release.Minor, ft)
	// Add Distro if there is one
	dist := ""
	if v.Distro != "" {
		dist = "-" + string(v.Distro)
	}
	// Add Enterprise modifier
	ent := ""
	if v.Release.Enterprise {
		ent = "-enterprise"
	}
	return "mongodb-" + string(v.OS) + "-" + string(v.Arch) + ent + dist + rel, nil
}

// Get directory for a Version
func (v *Version) ToDirectory() (string, error) {
	err := v.Validate()
	if err != nil {
		return "", err
	}
	if v.OS == "macos" {
		return "osx/", nil
	}
	return string(v.OS) + "/", nil
}

// Get URL for a Version
func (v *Version) ToURL() (string, error) {
	dir, err := v.ToDirectory()
	if err != nil {
		return "", err
	}
	fn, err := v.ToFilename()
	if err != nil {
		return "", err
	}
	prefix := communityUrlPrefix
	if v.Release.Enterprise {
		prefix = enterpriseUrlPrefix
	}
	return prefix + dir + fn, nil
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
	for _, x := range validOs {
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
