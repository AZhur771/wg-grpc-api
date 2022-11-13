package main

import (
	"flag"
	"fmt"

	"github.com/AZhur771/wg-grpc-api/internal/config"
)

var configFile string

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/wgapi/config.yaml", "Path to configuration file")
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Printf("Release: %s\nBuild date: %s\nGit hash: %s", release, buildDate, gitHash)
		return
	}

	config, err := config.ParseConfig(configFile)
	panicOnErr(err)

	_ = config
}
