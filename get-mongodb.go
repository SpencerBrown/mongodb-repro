package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/SpencerBrown/mongodb-repro/config"
	"github.com/SpencerBrown/mongodb-repro/get"
	"github.com/SpencerBrown/mongodb-repro/version"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const binaryDir = "mongodb-binaries"
const runtimeDir = "mongodb-runtime"

var binaryPath string
var runtimePath string

func init() {
	binaryPath = getPath(binaryDir)
	runtimePath = getPath(runtimeDir)
}

func getPath(dir string) string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v\n", err)
	}
	return filepath.Join(homedir, dir)
}

func main() {

	var arch = flag.String("arch", "x86_64", "Architecture: x86_64 (default), aarch64, ppc64ld, s390x")
	//var a_arch = flag.String("a", "x86_64", "(short for 'arch')")
	var myos = flag.String("os", "linux", "OS: linux (default), macos, win32")
	//var a_myos = flag.String("o", "linux", "(short for 'os')")
	var distro = flag.String("distro", "ubuntu1604", "Linux distro")
	var release = flag.String("version", "4.2.9", "MongoDB version")
	var community = flag.Bool("community", false, "Community version?")

	flag.Parse()

	if flag.NArg() != 1 {
		printHelp()
		return
	}
	switch flag.Arg(0) {
	case "marshal":
		cfg := &config.Type{}
		buf := cfg.ToYaml()
		_, _ = buf.WriteTo(os.Stdout)
	case "list":
		err := listVersions()
		if err != nil {
			fmt.Printf("Error listing versions: %v\n", err)
		}
	case "get":
		v := &version.Version{
			Arch:   version.ArchType(*arch),
			OS:     version.OSType(*myos),
			Distro: version.DistroType(*distro),
		}
		var err error
		v.Release, err = version.ToRelease(*release)
		if err != nil {
			fmt.Printf("Error in release '%s': %v\n", *release, err)
			return
		}
		v.Release.Enterprise = !(*community)
		fmt.Println(v)
		err = v.Validate()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		err = getOneAndExpand(v)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	case "config":
		//config.Foo()
		var cfg config.Type
		cfg = *config.OurDefaults // makes a copy so we don't pollute the static global variable. This makes a full copy because we don't have any reference types in the struct.
		cfg.Storage.DbPath = filepath.Join(runtimePath, "data")
		cfg.SystemLog.Path = filepath.Join(runtimePath, "sa.log")
		err := config.WriteConfig(&cfg, runtimePath, "sa.yaml")
		if err == nil {
			fmt.Printf("Configuration complete!\n")
			buf := cfg.ToYaml()
			_, _ = buf.WriteTo(os.Stdout)
		} else {
			fmt.Printf("Setup error: %v\n", err)
			return
		}
	case "run":
		mongoExt := ""
		if runtime.GOOS == "windows" {
			mongoExt = ".exe"
		}
		runcmd := exec.Command(filepath.Join(binaryPath, "mongodb-win32-x86_64-enterprise-windows-64-4.2.8/bin/mongod"+mongoExt), "-f", filepath.Join(runtimePath, "sa.yaml"))
		err := runcmd.Start()
		if err != nil {
			fmt.Printf("Error ) MongoDB: %v\n", err)
			return
		}
		fmt.Printf("Started!\n")
		client, err := connectMongo()
		if err != nil {
			fmt.Printf("Error connecting to server: %v\n", err)
			return
		}
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()
		err = setupAdminUser(client)
		if err != nil {
			fmt.Printf("Error setting up admin user: %v\n", err)
			return
		}
		fmt.Printf("Successfully set up admin user!\n")
	default:
		printHelp()
	}
}

func setupAdminUser(client *mongo.Client) error {
	cmd := bson.D{{"createUser", "admin"}, {"pwd", "tester"}, {"roles", bson.A{"root"}}}
	db := client.Database("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := db.RunCommand(ctx, cmd)
	if res.Err() != nil {
		return fmt.Errorf("error running createUser command: %v", res.Err())
	}
	var result bson.M
	_ = res.Decode(&result)
	fmt.Printf("Created asmin user: \n%v\n", result)
	return nil
}

func connectMongo() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("error setting up client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("error connecting: %v", err)
	}
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("error pinging: %v", err)
	}
	return client, nil
}

func printHelp() {
	fmt.Printf("%s list - lists currently downloaded versions\n", os.Args[0])
	fmt.Printf("%s get - downloads a version\n", os.Args[0])
	flag.PrintDefaults()
}

func getOneAndExpand(v *version.Version) error {
	myLocation, err := v.ToLocation()
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
	//myVersion2, err := version.ToVersion(myLocation.Filename)
	//fmt.Println(myVersion2)
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
		fmt.Printf("%-7s %-5s %-10s %d.%d.%2d %s %s\n", v.Arch, v.OS, v.Distro, v.Release.Version, v.Release.Major, v.Release.Minor, v.Release.Modifier, isEnterprise)
	}
	return nil
}
