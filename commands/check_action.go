package commands

import (
	"os"
	"path/filepath"
	"strings"
	"errors"
	"io/ioutil"

	"github.com/urfave/cli"
	"github.com/Droi-SDK/droi-checker/logger"
)

func checkAction(c *cli.Context) error {
	path := c.Args().First()

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
	logger.Info("Current work dir:" + absPath)
	platform := whichPlatform(absPath)

	switch platform {
	case "android":
		logger.Info("Platform:" + "android")
		androidChecker(absPath)
	case "ios":
		logger.Info("Platform:" + "ios")
	default:
		logger.Info("Platform:" + "unknown")
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
