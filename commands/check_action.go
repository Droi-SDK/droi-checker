package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"errors"
	"io/ioutil"

	"github.com/urfave/cli"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var(
	output = colorable.NewColorableStderr()
)

func checkAction(c *cli.Context) error {
	path := c.Args().First()
	prefix := color.GreenString("[work dir]")

	if strings.HasPrefix(path, "-") {
		return errors.New("Don't use - prefix")
	}
	if path == "" {
		path = getSrcFullPath()
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fmt.Fprintln(output, prefix, absPath)
	platform := whichPlatform(absPath)
	platformPrefix := color.GreenString("[platform]")

	switch platform {
	case "android":
		fmt.Fprintln(output, platformPrefix, "android")
		androidChecker(absPath)
	case "ios":
		fmt.Fprintln(output, platformPrefix, "ios")
	default:
		fmt.Fprintln(output, platformPrefix, "unknown")
		errors.New("Unknown platform")
	}
	return nil
}

func whichPlatform(fullPath string) (platform string) {
	dir, _ := ioutil.ReadDir(fullPath)
	for _, fi := range dir {
		fileName := fi.Name();
		if strings.EqualFold(fileName, "build.gradle") {
			platform = "android"
			return
		}
		if strings.HasSuffix(fi.Name(), ".xcodeproj") ||
			strings.HasSuffix(fileName, ".xcworkspace") {
			platform = "ios"
			return
		}
	}
	return "unknown"
}

func getSrcFullPath() (fullPath string) {
	args := os.Args
	parameterLen := len(args)
	if parameterLen == 1 {
		fullPath, _ = os.Getwd()
	}
	fullPath, _ = filepath.Abs(fullPath)
	return
}