package commands

import (
	"strings"
	"github.com/Droi-SDK/droi-checker/logger"
)

func depCheck(dep compile){
	switch dep.artifactId {
	case "Core":
		coreCheck(dep)
	case "push":
		pushCheck(dep)
	case "feedback":
		feedbackCheck(dep)
	case "selfupdate":
		selfupdateCheck(dep)
	case "analytics":
		analyticsCheck(dep)
	default:
	}
}

// TODO 各个sdk在java中的初始化

func coreCheck(compile compile) {
	versionCheck(compile)
}

func pushCheck(compile compile) {
	versionCheck(compile)
}

func selfupdateCheck(compile compile) {
	versionCheck(compile)
}

func analyticsCheck(compile compile) {
	versionCheck(compile)
}

func feedbackCheck(compile compile) {
	versionCheck(compile)
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
