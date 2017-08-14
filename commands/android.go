package commands

import (
	"os"
	"path/filepath"
	"fmt"
	"regexp"
	"io/ioutil"
	"github.com/Droi-SDK/droi-checker/logger"
	"strings"
)

func androidChecker(path string) error {
	//findMainAppDir(path)
	settingGradlePath := filepath.Join(path, "settings.gradle")
	settingGradleFile := loadFile(settingGradlePath)
	result := parseSetting(settingGradleFile)
	for i:= range result{
		modulePath := result[i][1]
		mainBuildGradlePath := filepath.Join(path, modulePath, "build.gradle")
		mainBuildGradleFile := loadFile(mainBuildGradlePath)
		isApplication := isApplication(mainBuildGradleFile)
		logger.Info(modulePath,isApplication)
		if  isApplication{
			parseManifest(modulePath)
			parseBuildGradle(modulePath)
		}
	}
	return nil
}

func findAllModule(path string) {
	settingPath := filepath.Join(path, "setting.gradle")
	settingFile, err := os.Open(settingPath)
	defer settingFile.Close()
	if err != nil {
		fmt.Println(settingFile, err)
		return
	}
	buf := make([]byte, 1024)
	for {
		n, _ := settingFile.Read(buf)
		if 0 == n {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}

func parseSetting(settingGradle []byte) ([][]string) {
	reg1 := regexp.MustCompile(`include[ ]*[\'\"]\:(\w+)[\'\"]`)
	dataSlice1 := reg1.FindAllSubmatchIndex(settingGradle, -1)
	result1 := make([]string, len(dataSlice1))
	for i := range dataSlice1 {
		result1[i] = string(settingGradle[dataSlice1[i][2]:dataSlice1[i][3]])
	}
	fmt.Println(result1)

	reg2 := regexp.MustCompile(`project\([\s]*[\'\"]\:(\w+)[\'\"][\s]*\)[\s]*\.[\s]*projectDir[ \t]*=[\s]*new[ ]+File[\s]*\([\s]*[\'\"]([\w\\\/]+)[\'\"][\s]*\)`)
	dataSlice2 := reg2.FindAllSubmatchIndex(settingGradle, -1)
	result2 := make([][]string, len(dataSlice2))
	for i := range dataSlice2 {
		result2[i] = make([]string, 2)
		result2[i][0] = string(settingGradle[dataSlice2[i][2]:dataSlice2[i][3]])
		result2[i][1] = string(settingGradle[dataSlice2[i][4]:dataSlice2[i][5]])
	}

	for i := range result1 {
		if inArray(result1[i], result2) {
			continue
		} else {
			x := make([]string, 2)
			x[0] = result1[i]
			x[1] = result1[i]
			result2 = append(result2, x)
		}
	}
	fmt.Println(result2)
	return result2
}

func inArray(name string, array [][]string) (b bool) {
	b = false
	for i := range array {
		if strings.EqualFold(name, array[i][0]) {
			b = true
			return
		}
	}
	return
}

func loadFile(path string) (content []byte) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}
	content, _ = ioutil.ReadAll(file)
	return
}

func isApplication(mainBuildGradle []byte) bool {
	pattern := `apply[ ]+plugin:[\s]+'com.android.application'`
	if ok, _ := regexp.Match(pattern, mainBuildGradle); ok {
		return true
	} else {
		return false
	}
}

func parseManifest(modulePath string) {

}

func parseBuildGradle(modulePath string) {

}
