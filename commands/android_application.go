package commands

import (
	"fmt"
	"path/filepath"
	"os"
	"regexp"
	"github.com/Droi-SDK/droi-checker/logger"
)

func parseApplication(applicationPath string) {
	applicationFile := loadFile(applicationPath)
	fmt.Println(string(applicationFile))
}

func checkInitialize(javaPath string, compileArray []compile) (bool) {
	files, _ := WalkDir(javaPath)
	for i := range files {
		b, setChannel := findInit(files[i], compileArray)
		if b {
			return setChannel
			break
		}
	}
	return false
}

func findInit(path string, compileArray []compile) (b bool, setChannel bool) {
	file := loadFile(path)
	regInit := regexp.MustCompile(`Core[\s]*.[\s]*initialize[\s]*\([\s]*this[\w\.]*[\s]*\)[\s]*;`)
	initSlice := regInit.FindAllSubmatchIndex(file, -1)
	if len(initSlice) != 0 {
		regSetChannelName := regexp.MustCompile(`Core[\s]*.[\s]*setChannelName[\s]*\([\s]*\"([\w]+)\"[\s]*\)[\s]*;`)
		channelNameSlice := regSetChannelName.FindAllSubmatchIndex(file, -1)
		if len(channelNameSlice) != 0 {
			setChannel = true
			if channelNameSlice[0][0] < initSlice[0][0] {
				logger.Info("Core.setChannel()使用正确,channel name is:",
					string(file[channelNameSlice[0][2]:channelNameSlice[0][3]]))
			} else {
				logger.Error("Core.setChannel()使用错误，请检查是否在Core.intialize()之前调用")
			}
		}
		b = true
		for i:= range compileArray {
			logger.Error("xxx")
			sdkName := converArtifactId(compileArray[i].artifactId)
			regSdk := regexp.MustCompile(sdkName+`[\s]*.[\s]*initialize[\s]*\([\s]*this[\w\.]*[\s]*[\,][\s]*[\w\_][\s]*\)[\s]*;`)
			sdkSlice := regSdk.FindAllSubmatchIndex(file, -1)
			if len(sdkSlice) != 0 {
				setChannel = true
				if sdkSlice[0][0] > initSlice[0][0] {
					logger.Info(sdkName+"initialize",
						string(file[sdkSlice[0][2]:sdkSlice[0][3]]))
				} else {
					logger.Error(sdkName+"initialize使用错误，请检查是否在Core.intialize()之后调用")
				}
			}
		}
	}
	return
}

func converArtifactId(artifactId string) (sdkName string) {
	switch artifactId {
	case "feedback":
		sdkName = "DroiFeedback"
	case "selfupdate":
		sdkName = "DroiUpdate"
	case "analytics":
		sdkName = "DroiAnalytics"
	case "push":
		sdkName = "DroiPush"
	}
	return
}

func WalkDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 30)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		files = append(files, filename)
		return nil
	})
	return files, err
}
