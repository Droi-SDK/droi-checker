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
	//getProjectExt(path)
	for i := range result {
		modulePath := result[i].path
		moduleAbsPath := filepath.Join(path, modulePath)
		mainBuildGradlePath := filepath.Join(path, modulePath, "build.gradle")
		mainBuildGradleFile := loadFile(mainBuildGradlePath)
		isApplication := isApplication(mainBuildGradleFile)
		logger.Info("module:", modulePath, "is Application")
		if isApplication {
			droibaasCompile := parseBuildGradle(moduleAbsPath)
			//getExt(modulePath, string(mainBuildGradleFile))
			depCheck(droibaasCompile)
			ok, _, setChannel1, channelName1 := parseManifest(moduleAbsPath)
			if !ok {
				logger.Error("解析manifest错误")
			}
			javaPath := filepath.Join(moduleAbsPath, "src", "main", "java")
			setChannel2, channelName2 := checkInitialize(javaPath, droibaasCompile)
			if !setChannel1 && !setChannel2 {
				logger.Warn("请设置渠道号，否则将使用默认渠道号：UNKNOWN_CHANNEL")
			} else if setChannel1 && !setChannel2 {
				logger.Info("channel name is；" + channelName1)
			} else if !setChannel1 && setChannel2 {
				logger.Info("channel name is；" + channelName2)
			} else if setChannel1 && setChannel2 {
				logger.Info("以代码设置为准，channel name is；" + channelName2)
			}
		}
	}
	return nil
}

//func getProjectExt(path string) (extMap map[string]string) {
//	extMap = make(map[string]string)
//	buildGradlePath := filepath.Join(path, "build.gradle")
//	buildGradleFile := loadFile(buildGradlePath)
//	subprojectsReg := regexp.MustCompile(`subprojects[\s]*{([\d\D]*)}`)
//	dataSlice1 := subprojectsReg.FindAllSubmatchIndex(buildGradleFile, -1)
//	var subprojectsExt map[string]string //1
//	if len(dataSlice1) != 0 {
//		extString := string(buildGradleFile[dataSlice1[0][0]:dataSlice1[0][1]])
//		subprojectsExt = getExt(extString)
//	}
//
//	allprojectsReg := regexp.MustCompile(`allprojects[\s]*{([\d\D]*)}`)
//	var allprojectsExt map[string]string //2
//	dataSlice2 := allprojectsReg.FindAllSubmatchIndex(buildGradleFile, -1)
//	if len(dataSlice2) != 0 {
//		extString := string(buildGradleFile[dataSlice1[0][0]:dataSlice1[0][1]])
//		allprojectsExt = getExt(extString)
//	}
//
//	projectExt := getExt(string(buildGradleFile))//3
//
//
//	buildScriptReg := regexp.MustCompile(`buildscript[\s]*{([\d\D]*)}`)
//	var buildScriptExt map[string]string //4
//	dataSlice3 := buildScriptReg.FindAllSubmatchIndex(buildGradleFile, -1)
//	if len(dataSlice3) != 0 {
//		extString := string(buildGradleFile[dataSlice1[0][0]:dataSlice1[0][1]])
//		allprojectsExt = getExt(extString)
//	}
//	return
//}

//func getExt(mainBuildGradle string) (extMap map[string]string) {
//	reg1 := regexp.MustCompile(`ext[\s]*{([\d\D]*)}`)
//	dataSlice1 := reg1.FindAllStringSubmatchIndex(mainBuildGradle, -1)
//	if len(dataSlice1) == 0 {
//		return
//	}
//	extMap = make(map[string]string)
//	for j := range dataSlice1 {
//		extString := string(mainBuildGradle[dataSlice1[j][2]:dataSlice1[j][3]])
//		regExt := regexp.MustCompile(`[ ]*([\w^=]*)[ ]*=[\s]*([\w\-.+"']*)`)
//		extSlice := regExt.FindAllStringSubmatchIndex(extString, -1)
//		logger.Error(len(extSlice))
//		for i := range extSlice {
//			extName := string(extString[extSlice[i][2]:extSlice[i][3]])
//			extValue := string(extString[extSlice[i][4]:extSlice[i][5]])
//			if strings.HasPrefix(extValue, "\"") == strings.HasSuffix(extValue, "\"") {
//				if strings.HasPrefix(extValue, "\"") {
//					extValue = strings.Trim(extValue, "\"")
//				}
//			}
//			extMap[extName] = extValue
//		}
//	}
//	return
//}

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
	reg1 := regexp.MustCompile(`include[ ]*['"]:(\w+)['"]`)
	dataSlice1 := reg1.FindAllSubmatchIndex(settingGradle, -1)
	result1 := make([]string, len(dataSlice1))
	for i := range dataSlice1 {
		result1[i] = string(settingGradle[dataSlice1[i][2]:dataSlice1[i][3]])
	}

	reg2 := regexp.MustCompile(`project\([\s]*['"]\:(\w+)['"][\s]*\)[\s]*\.[\s]*projectDir[ \t]*=[\s]*new[ ]+File[\s]*\([\s]*['"]([\w\\\/]+)['"][\s]*\)`)
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
	reg1 := regexp.MustCompile(`dependencies[\s]*{([\d\D]*)}`)
	dataSlice1 := reg1.FindAllSubmatchIndex(mainBuildGradleFile, -1)
	dependencies := string(mainBuildGradleFile[dataSlice1[0][2]:dataSlice1[0][3]])
	// 解析compile
	compileArray = make([]compile, 0)
	regCompile := regexp.MustCompile(`compile[ ]*[(]?[ ]*['"]([\w.\-]*):([\w.\-]*):([\w\-.@{}$+]*)['"][ ]*[)]?`)
	compileSlice := regCompile.FindAllStringSubmatchIndex(dependencies, -1)
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
	// 解析implementation
	regImplementation := regexp.MustCompile(`implementation[ ]*[(]?[ ]*['"]([\w.\-]*):([\w\.\-]*):([\w\-.@{}$+]*)['"][ ]*[\)]?`)
	implementationSlice := regImplementation.FindAllStringSubmatchIndex(dependencies, -1)

	for i := range implementationSlice {
		if len(implementationSlice[i]) != 8 {
			continue
		}
		groupId := dependencies[implementationSlice[i][2]:implementationSlice[i][3]]

		if !strings.EqualFold(groupId, "com.droi.sdk") {
			continue
		}
		artifactId := dependencies[implementationSlice[i][4]:implementationSlice[i][5]]
		version := dependencies[implementationSlice[i][6]:implementationSlice[i][7]]
		oneCompile := compile{groupId, artifactId, version}
		compileArray = append(compileArray, oneCompile)
	}
	// 解析config
	//regConfig := regexp.MustCompile(`defaultConfig[\s]*{[\s]*[\d\D]*?manifestPlaceholders[ ]*=[\s]*\[([\d\D]+?)][\s]*}[\s]*`)
	//configSlice := regConfig.FindAllSubmatchIndex(mainBuildGradleFile, -1)
	//var value []string = make([]string, len(configSlice))
	//for i := range configSlice {
	//	value[i]= string(mainBuildGradleFile[configSlice[i][2]:configSlice[i][3]])
	//}
	//logger.Error(value)
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
