package config

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sync"
)

var (
	attributeFlag = "#attr"  //root.ip#attr
	valueFlag     = "#value" //root.ip#value
)

type XmlConfig struct {
}

type XmlConfigManager struct {
	filename       string
	data           map[string]map[string]string // section=> key:val
	sectionComment map[string]string            // section : comment
	keyComment     map[string]string            // id: []{comment, key...}; id 1 is for main comment.
	changed        bool
	sync.RWMutex
}

func (xmlconf *XmlConfig) Parse(filename string) (ConfigManager, error) {
	return xmlconf.parseFile(filename)
}

func (xmlconf *XmlConfig) parseFile(filename string) (ConfigManager, error) {
	//文件是否存在
	xmlcontent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &IniConfigManager{filename: filename,
		data:           make(map[string]map[string]string),
		sectionComment: make(map[string]string),
		keyComment:     make(map[string]string),
		RWMutex:        sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()

	decoder := xml.NewDecoder(bytes.NewBuffer(xmlcontent))

	for element, err := decoder.Token(); err == nil; element, err = decoder.Token() {
		switch eType := element.(type) {
		case xml.StartElement:
			fmt.Println(eType)
		case xml.CharData:
		case xml.EndElement:
		case xml.Comment:
		}
	}
	return nil, nil

}
