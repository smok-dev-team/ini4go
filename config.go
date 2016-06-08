package config

import (
	"regexp"
	"io"
	"strings"
	"bufio"
	"os"
	"sync"
)

const (
	k_DEFAULT_SECTION = "default"
)

////////////////////////////////////////////////////////////////////////////////
var sectionRegexp = regexp.MustCompile(`\[(?P<header>[^]]+)\]$`)
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
	sections map[string]*Section
}

func (this *RawConfigParser) LoadFile(file string) error {
	var f, err = os.OpenFile("./test.conf", os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	return this.load(f)
}

func (this *RawConfigParser) load(r io.Reader) error {
	if this.sections == nil {
		this.sections = make(map[string]*Section)
	}

	var reader = bufio.NewReader(r)
	var line []byte
	var err error

	var currentSection *Section

	for {
		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		var sLine = strings.TrimSpace(string(line))

		// 如果是注释或者空行,则忽略
		if sLine == "" || strings.HasPrefix(sLine, "#") || strings.HasPrefix(sLine, ";") {
			continue
		}

		var sectionName = getSectionName(sLine)
		if len(sectionName) > 0 {
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

////////////////////////////////////////////////////////////////////////////////
func (this *RawConfigParser) NewSection(name string) *Section {
	this.Lock()
	defer this.Unlock()

	var section = this.sections[name]
	if section == nil {
		section = NewSection(name)
		this.sections[name] = section
	}
	return section
}

func (this *RawConfigParser) Section(name string) *Section {
	var s, _ = this.sections[name]
	return s
}

func (this *RawConfigParser) Sections() []string {
	var kList = make([]string, 0, len(this.sections))
	for key := range this.sections {
		kList = append(kList, key)
	}
	return kList
}

func (this *RawConfigParser) HasSection(section string) bool {
	var _, ok = this.sections[section]
	return ok
}

func (this *RawConfigParser) RemoveSection(section string) {
	delete(this.sections, section)
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
		return s.Options()
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
func (this *RawConfigParser) GetValue(section, option string) string {
	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.value[0]
		}
	}
	return ""
}

func (this *RawConfigParser) GetListValue(section, option string) []string {
	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.value
		}
	}
	return nil
}