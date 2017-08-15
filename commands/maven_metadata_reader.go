package commands

import (
	"encoding/xml"
	"github.com/levigross/grequests"
	"log"
)

type MavenMetadata struct {
	XMLName xml.Name `xml:"metadata"`
	Release string `xml:"versioning>release"`
}

func readFile(module string) (file []byte) {
	url := "https://raw.githubusercontent.com/DroiBaaS/DroiBaaS-SDK-Android/master/com/droi/sdk/" +
		module +
		"/maven-metadata.xml"
	resp, err := grequests.Get(url, nil)
	// You can modify the request by passing an optional RequestOptions struct

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}

	return resp.Bytes()
}

func parseMavenMetadata(module string) (lastestVersion string) {
	metaFile := readFile(module)
	var metadata MavenMetadata
	err := xml.Unmarshal(metaFile, &metadata)
	if err == nil {
		lastestVersion = metadata.Release
	}
	return
}
