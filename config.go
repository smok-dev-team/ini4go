package config

import (
	"regexp"
	"os"
)

////////////////////////////////////////////////////////////////////////////////
var sectionRegexp = regexp.MustCompile(`\[(?P<header>[^]]+)\]$`)
func sectionName(src string) (name string) {
	var rList = sectionRegexp.FindStringSubmatch(src)
	if len(rList) >= 2 {
		name = rList[1]
	}
	return name
}

var optionRegexp = regexp.MustCompile(`(?P<key>[^:=\s][^:=]*)\s*(?:(?P<vi>[:=])\s*(?P<value>.*))?$`)
func optionAndValue(src string) (option, vi, value string) {
	var rList = optionRegexp.FindStringSubmatch(src)
	if len(rList) >= 4 {
		option, vi, value = rList[1], rList[2], rList[3]
	}
	return option, vi, value
}

////////////////////////////////////////////////////////////////////////////////
type RawConfigParser struct {
	sections map[string]section
}

func (this *RawConfigParser) Load(file string) bool {

	return false
}

func (this *RawConfigParser) load(file os.File) {
}

func (this *RawConfigParser) Sections() []string {
	return nil
}