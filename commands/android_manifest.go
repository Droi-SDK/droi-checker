package commands

import (
	"path/filepath"
	"encoding/xml"
	"strings"
	"github.com/Droi-SDK/droi-checker/logger"
)

type Manifest struct {
	XMLName     xml.Name `xml:"manifest"`
	Application Application `xml:"application"`
	Package     string `xml:"package,attr"`
}

type Application struct {
	Metadatas []Metadata `xml:"meta-data"`
	Name      string `xml:"name,attr"`
}

type Metadata struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func parseManifest(modulePath string) (ok bool, applicationFilePath string, setChannel bool) {
	mainManifestPath := filepath.Join(modulePath, "src", "main", "AndroidManifest.xml")
	mainManifestFile := loadFile(mainManifestPath)
	//fmt.Println(string(mainManifestFile))
	var manifest Manifest
	err := xml.Unmarshal(mainManifestFile, &manifest)
	if err != nil {
		ok = false
	} else {
		ok = true
		if strings.HasPrefix(manifest.Application.Name, ".") {
			applicationFilePath = manifest.Package + manifest.Application.Name
		} else {
			applicationFilePath = manifest.Application.Name
		}
		setChannel = checkKeys(manifest)
	}
	return
}

func checkKeys(manifest Manifest) (bool) {
	metadatas := manifest.Application.Metadatas
	hasApplicationId, applicationIdIndex := hasKey(metadatas, "com.droi.sdk.application_id")
	if hasApplicationId {
		checkApplicationId(metadatas[applicationIdIndex].Value)
	} else {
		logger.Error("applicationId not set")
	}
	hasChannelName, channelNameIndex := hasKey(metadatas, "com.droi.sdk.channel_name")
	if hasChannelName {
		logger.Info("channel_name is:", metadatas[channelNameIndex].Value)
		return true
	} else {
		return false
	}
}

func checkApplicationId(applicationId string) {
	logger.Info("applicationId is:", applicationId)
}

func hasKey(metadatas []Metadata, key string) (b bool, index int) {
	b = false
	for i := range metadatas {
		if strings.EqualFold(metadatas[i].Name, key) {
			b = true
			index = i
			return
		}
	}
	return
}
