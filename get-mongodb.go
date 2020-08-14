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

	var arch = flag.String("arch", "x86_64", "Architecture: x86_64 (default), aarch64, ppc64le, s390x")
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
	err = v.Validate()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var isWindows = v.OS == "win32"

	switch flag.Arg(0) {
	case "marshal":
		fmt.Println("MongoDB Defaults applied")
		cfg := config.MongoDBDefaults
		buf := cfg.ToYaml(isWindows)
		_, _ = buf.WriteTo(os.Stdout)
		fmt.Println("Opinionated Defaults applied")
		cfg = config.OurDefaults
		buf = cfg.ToYaml(isWindows)
		_, _ = buf.WriteTo(os.Stdout)
		fmt.Println("Unmarshaling back to config struct")
		cfgbytes, err := ioutil.ReadFile(filepath.Join(runtimePath, "sa.yaml.gob"))
		if err != nil {
			fmt.Printf("Error reading GoB file %v\n", err)
			break
		}
		cfg2, err := config.FromGoB(cfgbytes)
		if err != nil {
			fmt.Printf("Error unmarshaling GoB file %v\n", err)
			break
		}
		fmt.Println(cfg2)
	case "list":
		err := listVersions()
		if err != nil {
			fmt.Printf("Error listing versions: %v\n", err)
		}
	case "get":
		err = getOneAndExpand(v)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "config":
		var cfg config.Type
		cfg = *config.OurDefaults // makes a copy so we don't pollute the static global variable. This makes a full copy because we don't have any reference types in the struct.
		cfg.Storage.DbPath = filepath.Join(runtimePath, "data")
		cfg.SystemLog.Path = filepath.Join(runtimePath, "sa.log")
		//cfg.ProcessManagement.Fork = true
		//isWindows = true
		err := config.WriteConfig(&cfg, runtimePath, "sa.yaml", isWindows)
		if err == nil {
			fmt.Printf("Configuration complete!\n")
		} else {
			fmt.Printf("Setup error: %v\n", err)
		}
	case "run":
		loc, err := v.ToLocation()
		if err != nil {
			fmt.Printf("Error converting to filename: %v\n", loc)
			break
		}
		fn := loc.Filename
		mongoExt := ""
		if isWindows {
			mongoExt = ".exe"
		}
		runcmd := exec.Command(filepath.Join(binaryPath, fn, "bin", "mongod"+mongoExt), "-f", filepath.Join(runtimePath, "sa.yaml"))
		err = runcmd.Start()
		if err != nil {
			fmt.Printf("Error ) MongoDB: %v\n", err)
			break
		}
		fmt.Printf("Started!\n")
		client, err := connectMongo(false)
		if err != nil {
			fmt.Printf("Error connecting to server: %v\n", err)
			break
		}
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()
		err = setupAdminUser(client)
		if err != nil {
			fmt.Printf("Error setting up admin user: %v\n", err)
			break
		}
		fmt.Printf("Successfully set up admin user!\n")
	case "stop":
		client, err := connectMongo(true)
		if err != nil {
			fmt.Printf("Error connecting to server: %v\n", err)
			break
		}
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()
		err = shutdownServer(client)
		if err != nil {
			fmt.Printf("Error shutting down server: %v\n", err)
			break
		}
		fmt.Printf("Successfully shut down server!\n")
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
	fmt.Printf("Created admin user: \n%v\n", result)
	return nil
}

func shutdownServer(client *mongo.Client) error {
	cmd := bson.D{{"shutdown", 1}}
	db := client.Database("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := db.RunCommand(ctx, cmd)
	if res.Err() != nil {
		return fmt.Errorf("error running shutdown command: %v", res.Err())
	}
	var result bson.M
	_ = res.Decode(&result)
	fmt.Printf("Shut down server: \n%v\n", result)
	return nil
}

func connectMongo(auth bool) (*mongo.Client, error) {
	copt := new(options.ClientOptions)
	copt.Hosts = []string{"localhost:27017"}
	if auth {
		copt.Auth = &options.Credential{
			Username: "admin",
			Password: "tester",
		}
	}
	client, err := mongo.NewClient(copt)
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
