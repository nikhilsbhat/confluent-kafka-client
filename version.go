package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	// Version specifies the version of the application and cannot be changed by end user.
	Version string

	// Env tells end user that what variant (here we use the name of the git branch to make it simple)
	// of application is he using.
	Env string

	BuildDate string
	GoVersion string
	Platform  string
	Revision  string
)

type BuildInfo struct {
	Version     string
	Revision    string
	Environment string
	BuildDate   string
	GoVersion   string
	Platform    string
}

// GetBuildInfo return the version and other build info of the application.
func GetBuildInfo() BuildInfo {
	if strings.ToLower(Env) != "production" {
		Env = "alfa"
	}
	return BuildInfo{
		Version:     Version,
		Revision:    Revision,
		Environment: Env,
		Platform:    Platform,
		BuildDate:   BuildDate,
		GoVersion:   GoVersion,
	}
}

func appVersion(c *cli.Context) error {
	buildInfo, err := json.Marshal(GetBuildInfo())
	if err != nil {
		return err
	}
	fmt.Println("renderer version:", string(buildInfo))
	return nil
}
