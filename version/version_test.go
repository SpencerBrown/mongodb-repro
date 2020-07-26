package version

import (
	"testing"
)

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
https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-4.2.8.tgz
https://fastdl.mongodb.org/osx/mongodb-osx-ssl-x86_64-4.0.19.tgz
https://downloads.mongodb.com/osx/mongodb-osx-x86_64-enterprise-4.0.19.tgz

*/

// mongodb-(linux|macos|osx|osx-ssl|win32)-(x86_64|s390x|ppc64le|aarch64)-(enterprise-)?(rhel62|rhel70|...|windows-64)?\d\.\d\.\d(-+*)?(.tgz|.zip)?

func TestVersion_ToLocation(t *testing.T) {
	type fields struct {
		arch    ArchType
		os      OSType
		distro  DistroType
		release ReleaseType
	}
	tests := []struct {
		name    string
		fields  fields
		want    Location
		wantErr bool
	}{
		{
			"windows",
			fields{"x86_64", "win32", "windows-64", ReleaseType{4, 2, 5, "", true}},
			Location{Filename: "mongodb-win32-x86_64-enterprise-windows-64-4.2.5", URLSuffix: ".zip", URLPrefix: "https://downloads.mongodb.com/win32/"},
			false,
		},
		{
			"mac",
			fields{"x86_64", "macos", "", ReleaseType{4, 2, 5, "", true}},
			Location{Filename: "mongodb-macos-x86_64-enterprise-4.2.5", URLPrefix: "https://downloads.mongodb.com/osx/", URLSuffix: ".tgz"},
			false,
		},
		{
			"mac-rc",
			fields{"x86_64", "macos", "", ReleaseType{4, 2, 5, "rc1", true}},
			Location{Filename: "mongodb-macos-x86_64-enterprise-4.2.5-rc1", URLPrefix: "https://downloads.mongodb.com/osx/", URLSuffix: ".tgz"},
			false,
		},
		{
			"mac40",
			fields{"x86_64", "macos", "", ReleaseType{4, 0, 19, "", true}},
			Location{Filename: "mongodb-osx-x86_64-enterprise-4.0.19", URLPrefix: "https://downloads.mongodb.com/osx/", URLSuffix: ".tgz"},
			false,
		},
		{
			"linux",
			fields{"s390x", "linux", "ubuntu1804", ReleaseType{4, 2, 5, "", true}},
			Location{Filename: "mongodb-linux-s390x-enterprise-ubuntu1804-4.2.5", URLPrefix: "https://downloads.mongodb.com/linux/", URLSuffix: ".tgz"},
			false,
		},
		{
			"linux/community",
			fields{"x86_64", "linux", "ubuntu1604", ReleaseType{4, 2, 5, "", false}},
			Location{Filename: "mongodb-linux-x86_64-ubuntu1604-4.2.5", URLPrefix: "https://fastdl.mongodb.org/linux/", URLSuffix: ".tgz"},
			false,
		},
		{
			"mac/community",
			fields{"x86_64", "macos", "", ReleaseType{4, 2, 8, "", false}},
			Location{Filename: "mongodb-macos-x86_64-4.2.8", URLPrefix: "https://fastdl.mongodb.org/osx/", URLSuffix: ".tgz"},
			false,
		},
		{
			"mac40/community",
			fields{"x86_64", "macos", "", ReleaseType{4, 0, 19, "", false}},
			Location{Filename: "mongodb-osx-ssl-x86_64-4.0.19", URLPrefix: "https://fastdl.mongodb.org/osx/", URLSuffix: ".tgz"},
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
			got, err := v.ToLocation()
			if err == nil {
				if tt.wantErr {
					t.Errorf("ToLocation(): Version %v: Wanted error, got none", v)
				} else {
					if (got == nil) || (got.Filename != tt.want.Filename) || (got.URLPrefix != tt.want.URLPrefix) {
						t.Errorf("ToLocation(): got %v, wanted %v", got, tt.want)
					}
				}
			} else {
				if !tt.wantErr {
					t.Errorf("ToLocation(): Version: %v, got unwanted error %v", v, err)
				}
			}
		})
	}
}

func TestVersion_ToVersion(t *testing.T) {
	tests := []struct {
		name    string
		fn      string
		want    Version
		wantErr bool
	}{
		{
			name:    "windows",
			fn:      "mongodb-win32-x86_64-enterprise-windows-64-4.2.5.zip",
			want:    Version{Arch: "x86_64", OS: "win32", Distro: "windows-64", Release: ReleaseType{Version: 4, Major: 2, Minor: 5, Modifier: "", Enterprise: true}},
			wantErr: false,
		}, {
			name:    "windows-rc",
			fn:      "mongodb-win32-x86_64-enterprise-windows-64-4.2.5-rc0.zip",
			want:    Version{Arch: "x86_64", OS: "win32", Distro: "windows-64", Release: ReleaseType{Version: 4, Major: 2, Minor: 5, Modifier: "rc0", Enterprise: true}},
			wantErr: false,
		},
		{
			name:    "mac",
			fn:      "mongodb-macos-x86_64-enterprise-4.2.5",
			want:    Version{"x86_64", "macos", "", ReleaseType{4, 2, 5, "", true}},
			wantErr: false,
		},
		{
			name:    "mac40",
			fn:      "mongodb-osx-x86_64-enterprise-4.0.19.tgz",
			want:    Version{"x86_64", "macos", "", ReleaseType{4, 0, 19, "", true}},
			wantErr: false,
		},
		{
			name:    "linux",
			fn:      "mongodb-linux-s390x-enterprise-ubuntu1804-4.2.5.tgz",
			want:    Version{"s390x", "linux", "ubuntu1804", ReleaseType{4, 2, 5, "", true}},
			wantErr: false,
		},
		{
			name:    "linux/community",
			fn:      "mongodb-linux-x86_64-ubuntu1604-4.2.5",
			want:    Version{"x86_64", "linux", "ubuntu1604", ReleaseType{4, 2, 5, "", false}},
			wantErr: false,
		},
		{
			name:    "mac/community",
			want:    Version{"x86_64", "macos", "", ReleaseType{4, 2, 8, "", false}},
			fn:      "mongodb-macos-x86_64-4.2.8.tgz",
			wantErr: false,
		},
		{
			name:    "mac40/community",
			want:    Version{"x86_64", "macos", "", ReleaseType{4, 0, 19, "", false}},
			fn:      "mongodb-osx-ssl-x86_64-4.0.19.tgz",
			wantErr: false,
		},
		{
			name:    "nomongo",
			fn:      "mangodb-foo-bar-what-ever",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badArch",
			fn:      "mongodb-macos-x86_65-4.2.8.tgz",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badOS",
			fn:      "mongodb-macosh-x86_64-4.2.8.tgz",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badenterprise",
			fn:      "mongodb-linux-s390x-enterprize-ubuntu1804-4.2.5.tgz",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badwindistro",
			fn:      "mongodb-win32-x86_64-enterprise-windoes-64-4.2.5.zip",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badmajor",
			fn:      "mongodb-win32-x86_64-enterprise-windows-64-4.two.5.zip",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badzip",
			fn:      "mongodb-win32-x86_64-enterprise-windows-64-4.2.5.tgz",
			want:    Version{},
			wantErr: true,
		},
		{
			name:    "badtgz",
			fn:      "mongodb-linux-x86_64-enterprise-amzn64-3.6.17.zip",
			want:    Version{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToVersion(tt.fn)
			//fmt.Println(got, err)
			if err == nil {
				if tt.wantErr {
					t.Errorf("ToVersion(): fn: %v: wanted error, got none", tt.fn)
				} else {
					if got == nil || got.OS != tt.want.OS || got.Arch != tt.want.Arch || got.Distro != tt.want.Distro || got.Release.Version != tt.want.Release.Version || got.Release.Major != tt.want.Release.Major || got.Release.Minor != tt.want.Release.Minor || got.Release.Enterprise != tt.want.Release.Enterprise {
						t.Errorf("ToVersion(): got %v, wanted %v", got, tt.want)
					}
				}
			} else {
				if !tt.wantErr {
					t.Errorf("ToVersion(): fn '%v': got unwanted error %v", tt.fn, err)
				}
			}
		})
	}
}

func TestVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       Version
		wantErr bool
	}{
		{
			name: "good",
			v: Version{
				Arch:   "x86_64",
				OS:     "linux",
				Distro: "amazon2",
				Release: ReleaseType{
					Version:    3,
					Major:      6,
					Minor:      12,
					Enterprise: false,
				},
			},
			wantErr: false,
		},
		{
			"bad Arch",
			Version{"foobar", "linux", "ubuntu1804", ReleaseType{4, 2, 5, "", true}},
			true,
		},
		{
			"bad OS",
			Version{"s390x", "foobar", "ubuntu1804", ReleaseType{4, 2, 5, "", true}},
			true,
		},
		{
			"bad Distro",
			Version{"s390x", "linux", "foobar", ReleaseType{4, 2, 5, "", true}},
			true,
		},
		{
			"bad Version",
			Version{"s390x", "linux", "ubuntu1804", ReleaseType{42, 2, 5, "", true}},
			true,
		},
		{
			"bad Major",
			Version{"s390x", "linux", "ubuntu1804", ReleaseType{4, 42, 5, "", true}},
			true,
		},
		{
			"bad minor",
			Version{"s390x", "linux", "ubuntu1804", ReleaseType{4, 2, 42, "", true}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if err == nil {
				if tt.wantErr {
					t.Errorf("Validate(): wanted error, got none: %v", tt.v)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("Validate(): got unwanted error %v on %v", err, tt.v)
				}
			}
		})
	}
}
