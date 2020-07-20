package version

import "testing"

/*
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-amazon2-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-amzn64-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-debian92-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-debian10-4.2.5.tgz
https://downloads.mongodb.com/osx/mongodb-macos-x86_64-enterprise-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-rhel62-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-s390x-enterprise-rhel67-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-rhel70-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-ppc64le-enterprise-rhel71-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-s390x-enterprise-rhel72-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-rhel80-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-s390x-enterprise-suse12-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-suse12-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-suse15-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-aarch64-enterprise-ubuntu1604-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-ubuntu1604-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-ppc64le-enterprise-ubuntu1804-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-s390x-enterprise-ubuntu1804-4.2.5.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-ubuntu1804-4.2.5.tgz
https://downloads.mongodb.com/win32/mongodb-win32-x86_64-enterprise-windows-64-4.2.5.zip

https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-amazon2-4.0.17.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-amzn64-3.6.17.tgz
https://downloads.mongodb.com/linux/mongodb-linux-x86_64-enterprise-amzn64-3.4.24.tgz

https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1604-4.2.5.tgz

*/

func TestVersion_ToDirectory(t *testing.T) {
	type fields struct {
		arch    ArchType
		os      OSType
		distro  DistroType
		release ReleaseType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"windows",
			fields{"x86_64", "win32", "windows-64", ReleaseType{4, 2, 0, true}},
			"win32/",
			false,
		},
		{
			"mac",
			fields{"x86_64", "macos", "", ReleaseType{4, 2, 0, true}},
			"osx/",
			false,
		},
		{
			"linux",
			fields{"x86_64", "linux", "ubuntu1804", ReleaseType{4, 2, 0, true}},
			"linux/",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				Arch:    tt.fields.arch,
				OS:      tt.fields.os,
				Distro:  tt.fields.distro,
				Release: tt.fields.release,
			}
			got, err := v.ToDirectory()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_ToFilename(t *testing.T) {
	type fields struct {
		arch    ArchType
		os      OSType
		distro  DistroType
		release ReleaseType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"windows",
			fields{"x86_64", "win32", "windows-64", ReleaseType{4, 2, 5, true}},
			"mongodb-win32-x86_64-enterprise-windows-64-4.2.5.zip",
			false,
		},
		{
			"mac",
			fields{"x86_64", "macos", "", ReleaseType{4, 2, 5, true}},
			"mongodb-macos-x86_64-enterprise-4.2.5.tgz",
			false,
		},
		{
			"linux",
			fields{"s390x", "linux", "ubuntu1804", ReleaseType{4, 2, 5, true}},
			"mongodb-linux-s390x-enterprise-ubuntu1804-4.2.5.tgz",
			false,
		},
		{
			"linux/community",
			fields{"x86_64", "linux", "ubuntu1604", ReleaseType{4, 2, 5, false}},
			"mongodb-linux-x86_64-ubuntu1604-4.2.5.tgz",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				Arch:    tt.fields.arch,
				OS:      tt.fields.os,
				Distro:  tt.fields.distro,
				Release: tt.fields.release,
			}
			got, err := v.ToFilename()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFilename() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Validate(t *testing.T) {
	type fields struct {
		arch    ArchType
		os      OSType
		distro  DistroType
		release ReleaseType
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"bad Arch",
			fields{"foobar", "linux", "ubuntu1804", ReleaseType{4, 2, 5, true}},
			true,
		},
		{
			"bad OS",
			fields{"s390x", "foobar", "ubuntu1804", ReleaseType{4, 2, 5, true}},
			true,
		},
		{
			"bad Distro",
			fields{"s390x", "linux", "foobar", ReleaseType{4, 2, 5, true}},
			true,
		},
		{
			"bad Version",
			fields{"s390x", "linux", "ubuntu1804", ReleaseType{42, 2, 5, true}},
			true,
		},
		{
			"bad Major",
			fields{"s390x", "linux", "ubuntu1804", ReleaseType{4, 42, 5, true}},
			true,
		},
		{
			"bad minor",
			fields{"s390x", "linux", "ubuntu1804", ReleaseType{4, 2, 42, true}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				Arch:    tt.fields.arch,
				OS:      tt.fields.os,
				Distro:  tt.fields.distro,
				Release: tt.fields.release,
			}
			if err := v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(err)
			}
		})
	}
}

func TestVersion_toURL(t *testing.T) {
	type fields struct {
		arch    ArchType
		os      OSType
		distro  DistroType
		release ReleaseType
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"windows",
			fields{"x86_64", "win32", "windows-64", ReleaseType{4, 2, 5, true}},
			"https://downloads.mongodb.com/win32/mongodb-win32-x86_64-enterprise-windows-64-4.2.5.zip",
			false,
		},
		{
			"mac",
			fields{"x86_64", "macos", "", ReleaseType{4, 2, 5, true}},
			"https://downloads.mongodb.com/osx/mongodb-macos-x86_64-enterprise-4.2.5.tgz",
			false,
		},
		{
			"linux",
			fields{"s390x", "linux", "ubuntu1804", ReleaseType{4, 2, 5, true}},
			"https://downloads.mongodb.com/linux/mongodb-linux-s390x-enterprise-ubuntu1804-4.2.5.tgz",
			false,
		},
		{
			"linux/community",
			fields{"x86_64", "linux", "ubuntu1604", ReleaseType{4, 2, 5, false}},
			"https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1604-4.2.5.tgz",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				Arch:    tt.fields.arch,
				OS:      tt.fields.os,
				Distro:  tt.fields.distro,
				Release: tt.fields.release,
			}
			got, err := v.ToURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
