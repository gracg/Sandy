package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var VERSION string = "1.0"

type Config struct {
	Server ConfigServer
	Client ConfigClient
}

type ConfigServer struct {
	Port     int
	SqLiteDB string
}

type ConfigClient struct {
	Target            string
	AuthorizedKeyFile string
}

func main() {
	ServerFlagPtr := flag.Bool("server", false, "sets sandy to run in server mode")
	ClientFlagPtr := flag.Bool("client", false, "sets sandy to run in client mode")
	VersionFlagPtr := flag.Bool("version", false, "displays Sandy version.")

	flag.Parse()

	if *VersionFlagPtr == true {
		fmt.Println("Sandy version " + VERSION)
	}

	BinaryPath := BinaryPath()
	ConfigFileExist, err := exists(BinaryPath + "/config.json")

	if err != nil || ConfigFileExist == false {
		e := errors.New("Sandy was unable to find 'config.json' file.")
		panic(e)
	}

	var ConfigFile Config
	ConfigBytes, err := ioutil.ReadFile(BinaryPath + "/config.json")

	if err != nil {
		e := errors.New("Sandy was unable to read 'config.json' file.")
		panic(e)
	}

	err = json.Unmarshal(ConfigBytes, &ConfigFile)
	if err != nil {
		panic(err)
	}

	if *ServerFlagPtr == true {
		ServerMain(ConfigFile)
	} else if *ClientFlagPtr == true {
		ClientMain(ConfigFile)
	}
}

func BinaryPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
