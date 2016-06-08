package config

import (
	"regexp"
	"io"
	"strings"
	"bufio"
	"os"
	"sync"
	"github.com/smartwalle/going/container"
	"bytes"
)

const (
	k_DEFAULT_SECTION = "default"
)

////////////////////////////////////////////////////////////////////////////////
var sectionRegexp = regexp.MustCompile(`^\[(?P<header>[^]]+)\]$`)
func getSectionName(src string) (name string) {
	var rList = sectionRegexp.FindStringSubmatch(src)
	if len(rList) >= 2 {
		name = rList[1]
	}
	return name
}

var optionRegexp = regexp.MustCompile(`(?P<key>[^:=\s][^:=]*)\s*(?:(?P<vi>[:=])\s*(?P<value>.*))?$`)
func getOptionAndValue(src string) (option, vi, value string) {
	var rList = optionRegexp.FindStringSubmatch(src)
	if len(rList) >= 4 {
		option, vi, value = rList[1], rList[2], rList[3]
	}
	return option, vi, value
}

////////////////////////////////////////////////////////////////////////////////
type RawConfigParser struct {
	sync.RWMutex
	sectionKeys []string
	sections    map[string]*Section
}

func (this *RawConfigParser) LoadFiles(files ...string) error {
	for _, file := range files {
		var f, err = os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		err = this.load(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *RawConfigParser) load(r io.Reader) error {
	if this.sections == nil {
		this.sections = make(map[string]*Section)
	}

	var reader = bufio.NewReader(r)
	var line []byte
	var err error

	var currentSection *Section

	var index = 0
	for {
		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if index == 0 {
			line = bytes.TrimPrefix(line, []byte("\xef\xbb\xbf"))
		}
		index++

		var sLine = strings.TrimSpace(string(line))

		// 如果是注释或者空行,则忽略
		if sLine == "" || strings.HasPrefix(sLine, "#") || strings.HasPrefix(sLine, ";") {
			continue
		}

		var sectionName = getSectionName(sLine)
		if len(sectionName) > 0 && strings.ToLower(sectionName) != k_DEFAULT_SECTION {
			currentSection = this.NewSection(sectionName)
			continue
		}

		if currentSection == nil {
			currentSection = this.NewSection(strings.ToLower(k_DEFAULT_SECTION))
		}

		var optName, optIV, optValue = getOptionAndValue(sLine)
		optName = strings.TrimSpace(optName)
		optIV = strings.TrimSpace(optIV)
		optValue = strings.TrimSpace(optValue)

		if optName != "" {
			currentSection.NewOption(optName, optIV, optValue)
		}
	}

	return nil
}

// TODO 写入文件
func (this *RawConfigParser) WriteToFile(file string) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func (this *RawConfigParser) NewSection(name string) *Section {
	this.Lock()
	defer this.Unlock()

	var section = this.sections[name]
	if section == nil {
		section = NewSection(name)
		this.sections[name] = section
		this.sectionKeys = append(this.sectionKeys, name)
	}
	return section
}

func (this *RawConfigParser) Section(name string) *Section {
	var s, _ = this.sections[name]
	return s
}

func (this *RawConfigParser) SectionNames() []string {
	return this.sectionKeys
}

func (this *RawConfigParser) SectionList() []*Section {
	var sList = make([]*Section, 0, len(this.sections))
	for _, value := range this.sections {
		sList = append(sList, value)
	}
	return sList
}

func (this *RawConfigParser) HasSection(section string) bool {
	var _, ok = this.sections[section]
	return ok
}

func (this *RawConfigParser) RemoveSection(section string) {
	delete(this.sections, section)
	container.Remove(&this.sectionKeys, section)
}

////////////////////////////////////////////////////////////////////////////////
func (this *RawConfigParser) Option(section, option string) *Option {
	var s = this.Section(section)
	if s != nil {
		return s.Option(option)
	}
	return nil
}

func (this *RawConfigParser) Options(section string) []string {
	var s = this.Section(section)
	if s != nil {
		return s.OptionKeys()
	}
	return nil
}

func (this *RawConfigParser) OptionList(section string) []*Option {
	var s = this.Section(section)
	if s != nil {
		return s.OptionList()
	}
	return nil
}

func (this *RawConfigParser) HasOption(section, option string) bool {
	if s, ok := this.sections[section]; ok {
		if _, ok := s.options[option]; ok {
			return true
		}
	}
	return false
}

func (this *RawConfigParser) RemoveOption(section, option string) {
	this.Lock()
	defer this.Unlock()

	var s = this.sections[section]
	if s != nil {
		s.RemoveOption(option)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (this *RawConfigParser) SetValue(section, option string, value string) {
	var s = this.NewSection(section)
	s.NewOption(option, "=", value)
}

////////////////////////////////////////////////////////////////////////////////
func (this *RawConfigParser) GetValue(section, option string) string {
	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.Value()
		}
	}
	return ""
}

func (this *RawConfigParser) GetListValue(section, option string) []string {
	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.ListValue()
		}
	}
	return nil
}