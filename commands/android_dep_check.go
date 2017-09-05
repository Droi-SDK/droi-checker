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
	lastVersion := getLatestVersion(compile)
	checkResult := isValidVersion(version, lastVersion)
	if checkResult == 0 {
		logger.Info(compile.artifactId +
			" SDK Version is correct!")
	} else if checkResult == 1 {
		logger.Warn(compile.artifactId +
			" SDK Version is " + version + ",没有检查该变量的值，请暂时自行检查！最新版本是：" +
			lastVersion)
	} else if checkResult == 2 {
		logger.Warn(compile.artifactId +
			" SDK Version is not lastest,the lastest Version is:" +
			lastVersion)
	}
}

func getLatestVersion(compile compile) (lastVersion string) {
	return parseMavenMetadata(compile.artifactId)
}

func isValidVersion(version string, lastsVersion string) (int) {
	versionSplit := strings.Split(lastsVersion, ".")
	validVersion := make([]string, 4)
	validVersion = append(validVersion, "+")
	validVersion = append(validVersion, versionSplit[0]+".+")
	validVersion = append(validVersion, versionSplit[0]+versionSplit[1]+".+")
	validVersion = append(validVersion, lastsVersion)
	for i := range validVersion {
		if strings.EqualFold(version, validVersion[i]) {
			return 0
		} else if strings.HasPrefix(version, "$") || (strings.HasPrefix(version, "${") && strings.HasSuffix(version, "}")) {
			return 1
		}
	}
	return 2
}
