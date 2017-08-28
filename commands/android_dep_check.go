package commands

import (
	"strings"
	"github.com/Droi-SDK/droi-checker/logger"
)

func depCheck(compileArray []compile) {
	for i := range compileArray {
		dep := compileArray[i]
		versionCheck(dep);
	}
}

func versionCheck(compile compile) {
	version := compile.version
	// TODO ext
	lastVersion := getLatsetVersion(compile)
	isValidVersion := isValidVersion(version, lastVersion)
	if isValidVersion {
		logger.Info(compile.artifactId +
			" SDK Version is correct!")
	} else {
		logger.Warn(compile.artifactId +
			" SDK Version is not lastest,the lastest Version is:" +
			lastVersion)
	}
}

func getLatsetVersion(compile compile) (lastVersion string) {
	return parseMavenMetadata(compile.artifactId)
}

func isValidVersion(version string, lastsVersion string) (bool) {
	versionSplit := strings.Split(lastsVersion, ".")
	validVersion := make([]string, 4)
	validVersion = append(validVersion, "+")
	validVersion = append(validVersion, versionSplit[0]+".+")
	validVersion = append(validVersion, versionSplit[0]+versionSplit[1]+".+")
	validVersion = append(validVersion, lastsVersion)
	for i := range validVersion {
		if strings.EqualFold(version, validVersion[i]) {
			return true
		}
	}
	return false
}
