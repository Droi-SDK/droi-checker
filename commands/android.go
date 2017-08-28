package commands

import (
	"os"
	"path/filepath"
	"regexp"
	"io/ioutil"
	"github.com/Droi-SDK/droi-checker/logger"
	"strings"
)

func androidChecker(path string) error {
	settingGradlePath := filepath.Join(path, "settings.gradle")
	settingGradleFile := loadFile(settingGradlePath)
	result := parseSetting(settingGradleFile)
	for i := range result {
		modulePath := result[i].path
		moduleAbsPath := filepath.Join(path, modulePath)
		mainBuildGradlePath := filepath.Join(path, modulePath, "build.gradle")
		mainBuildGradleFile := loadFile(mainBuildGradlePath)
		isApplication := isApplication(mainBuildGradleFile)
		logger.Info("module:", modulePath, "is Application")
		getProjectExt()
		if isApplication {
			droibaasCompile := parseBuildGradle(moduleAbsPath)
			depCheck(droibaasCompile)
			// TODO
			/*applicationPathDot := parseManifest(moduleAbsPath)
			applicationPath := getApplicationPathFromDot(moduleAbsPath, applicationPathDot)
			parseApplication(applicationPath)*/
			ok, _, setChannel1 := parseManifest(moduleAbsPath)
			if !ok {
				logger.Error("解析manifest错误")
			}
			javaPath := filepath.Join(moduleAbsPath, "src", "main", "java")
			setChannel2 := checkInitialize(javaPath,droibaasCompile)
			if !setChannel1 && !setChannel2 {
				logger.Warn("请设置渠道号，否则将使用默认渠道号：UNKNOWN_CHANNEL")
			}
		}
	}
	return nil
}

func getProjectExt()  {
	// TODO ext
}

func getApplicationPathFromDot(moduleAbsPath string, applicationPathDot string) (applicationPath string) {
	paths := strings.Split(applicationPathDot, ".")
	paths[len(paths)-1] += ".java"
	applicationPath = filepath.Join(moduleAbsPath, "src", "main", "java")
	for i := range paths {
		applicationPath = filepath.Join(applicationPath, paths[i])
	}
	return
}

func parseSetting(settingGradle []byte) ([]includeProject) {
	reg1 := regexp.MustCompile(`include[ ]*[\'\"]\:(\w+)[\'\"]`)
	dataSlice1 := reg1.FindAllSubmatchIndex(settingGradle, -1)
	result1 := make([]string, len(dataSlice1))
	for i := range dataSlice1 {
		result1[i] = string(settingGradle[dataSlice1[i][2]:dataSlice1[i][3]])
	}

	reg2 := regexp.MustCompile(`project\([\s]*[\'\"]\:(\w+)[\'\"][\s]*\)[\s]*\.[\s]*projectDir[ \t]*=[\s]*new[ ]+File[\s]*\([\s]*[\'\"]([\w\\\/]+)[\'\"][\s]*\)`)
	dataSlice2 := reg2.FindAllSubmatchIndex(settingGradle, -1)
	result2 := make([]includeProject, len(dataSlice2))
	for i := range dataSlice2 {
		result2[i].name = string(settingGradle[dataSlice2[i][2]:dataSlice2[i][3]])
		result2[i].path = string(settingGradle[dataSlice2[i][4]:dataSlice2[i][5]])
	}

	for i := range result1 {
		if inArray(result1[i], result2) {
			continue
		} else {
			include := includeProject{result1[i], result1[i]}
			result2 = append(result2, include)
		}
	}
	return result2
}

func inArray(name string, array []includeProject) (b bool) {
	b = false
	for i := range array {
		if strings.EqualFold(name, array[i].name) {
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

func parseBuildGradle(modulePath string) (compileArray []compile) {
	mainBuildGradlePath := filepath.Join(modulePath, "build.gradle")
	mainBuildGradleFile := loadFile(mainBuildGradlePath)
	// 找出dependencies部分
	reg1 := regexp.MustCompile(`dependencies[\s]*\{([\d\D]*)\}`)
	dataSlice1 := reg1.FindAllSubmatchIndex(mainBuildGradleFile, -1)
	dependencies := string(mainBuildGradleFile[dataSlice1[0][2]:dataSlice1[0][3]])
	// 解析compile
	regCompile := regexp.MustCompile(`compile[ ]*[\(]?[ ]*[\'\"]([\w\.\-]*)\:([\w\.\-]*)\:([\w\.\-\@\{\}\$\+]*)[\'\"][ ]*[\)]?`)
	compileSlice := regCompile.FindAllStringSubmatchIndex(dependencies, -1)
	compileArray = make([]compile, 0)
	for i := range compileSlice {
		if len(compileSlice[i]) != 8 {
			continue
		}
		groupId := dependencies[compileSlice[i][2]:compileSlice[i][3]]
		if !strings.EqualFold(groupId, "com.droi.sdk") {
			continue
		}
		artifactId := dependencies[compileSlice[i][4]:compileSlice[i][5]]
		version := dependencies[compileSlice[i][6]:compileSlice[i][7]]
		oneCompile := compile{groupId, artifactId, version}
		compileArray = append(compileArray, oneCompile)
	}
	// 解析config
	/*regConfig := regexp.MustCompile(`defaultConfig[\s]+\{[\s]+applicationId[ ]+[\'\"]([\w\.]*)[\'\"]`)
	configSlice := regConfig.FindAllSubmatchIndex(mainBuildGradleFile, -1)
	for i := range configSlice {
		value := string(mainBuildGradleFile[configSlice[i][2]:configSlice[i][3]])
		fmt.Println(value)
	}*/
	// TODO manifestPlaceholders 解析
	return
}

type compile struct {
	groupId    string
	artifactId string
	version    string
}

type includeProject struct {
	name string
	path string
}
