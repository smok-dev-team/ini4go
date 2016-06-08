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
	"fmt"
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
type rawConfigParser struct {
	sync.RWMutex
	sectionKeys []string
	sections    map[string]*Section
}

////////////////////////////////////////////////////////////////////////////////
type Config struct {
	rawConfigParser
}

func NewConfig() *Config {
	var c = &Config{}
	c.rawConfigParser.sections = make(map[string]*Section)
	return c
}

func (this *rawConfigParser) LoadFiles(files ...string) error {
	for _, file := range files {
		var f, err = os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		err = this.load(f)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *rawConfigParser) load(r io.Reader) error {
	var reader = bufio.NewReader(r)
	var line []byte
	var err error

	var currentSection *Section

	var index = 0
	var comments []string
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
		if sLine == "" {
			continue
		}

		if strings.HasPrefix(sLine, "#") || strings.HasPrefix(sLine, ";") {
			comments = append(comments, strings.TrimSpace(sLine[1:]))
			continue
		}

		var sectionName = getSectionName(sLine)
		if len(sectionName) > 0 && strings.ToLower(sectionName) != k_DEFAULT_SECTION {
			currentSection = this.NewSection(sectionName)
			currentSection.comments = append(currentSection.comments, comments...)
			comments = nil
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
			currentSection.NewOption(optName, optIV, optValue, comments)
			comments = nil
		}
	}
	return nil
}

func (this *rawConfigParser) WriteToFile(file string) error {
	var f, err = os.OpenFile(file, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return this.writeTo(f)
}

func (this *rawConfigParser) writeTo(w io.Writer) error {
	var writer = bufio.NewWriter(w)

	for _, sectionName := range this.sectionKeys {

		var section = this.Section(sectionName)
		for _, c := range section.Comments() {
			if len(strings.TrimSpace(c)) > 0 {
				writer.WriteString(fmt.Sprintf("# %s\n", c))
			}
		}

		writer.WriteString(fmt.Sprintf("[%s]\n", sectionName))

		for _, optionKey := range section.optionKeys {
			var opt = section.Option(optionKey)
			if len(opt.Comments()) > 0 {
				writer.WriteString("\n")
			}
			for _, c := range opt.Comments() {
				if len(strings.TrimSpace(c)) > 0 {
					writer.WriteString(fmt.Sprintf("# %s\n", c))
				}
			}
			for _, value := range opt.value {
				writer.WriteString(fmt.Sprintf("%s %s %s\n", opt.key, opt.iv, value))
			}
		}
		writer.WriteByte('\n')
	}
	return writer.Flush()
}

////////////////////////////////////////////////////////////////////////////////
func (this *rawConfigParser) NewSection(name string) *Section {
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

func (this *rawConfigParser) Section(name string) *Section {
	this.RLock()
	defer this.RUnlock()
	var s, _ = this.sections[name]
	return s
}

func (this *rawConfigParser) SectionNames() []string {
	return this.sectionKeys
}

func (this *rawConfigParser) SectionList() []*Section {
	var sList = make([]*Section, 0, len(this.sections))
	for _, value := range this.sections {
		sList = append(sList, value)
	}
	return sList
}

func (this *rawConfigParser) HasSection(section string) bool {
	var _, ok = this.sections[section]
	return ok
}

func (this *rawConfigParser) RemoveSection(section string) {
	this.Lock()
	defer this.Unlock()

	if strings.ToLower(section) == k_DEFAULT_SECTION {
		return
	}
	delete(this.sections, section)
	container.Remove(&this.sectionKeys, section)
}

////////////////////////////////////////////////////////////////////////////////
func (this *rawConfigParser) Option(section, option string) *Option {
	var s = this.Section(section)
	if s != nil {
		return s.Option(option)
	}
	return nil
}

func (this *rawConfigParser) Options(section string) []string {
	var s = this.Section(section)
	if s != nil {
		return s.OptionKeys()
	}
	return nil
}

func (this *rawConfigParser) OptionList(section string) []*Option {
	var s = this.Section(section)
	if s != nil {
		return s.OptionList()
	}
	return nil
}

func (this *rawConfigParser) HasOption(section, option string) bool {
	if s, ok := this.sections[section]; ok {
		if _, ok := s.options[option]; ok {
			return true
		}
	}
	return false
}

func (this *rawConfigParser) RemoveOption(section, option string) {
	this.Lock()
	defer this.Unlock()

	var s = this.sections[section]
	if s != nil {
		s.RemoveOption(option)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (this *rawConfigParser) SetValue(section, option string, value string) {
	var s = this.NewSection(section)
	s.NewOption(option, "=", value, nil)
}

////////////////////////////////////////////////////////////////////////////////
func (this *rawConfigParser) GetValue(section, option string) string {
	this.RLock()
	defer this.RUnlock()

	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.Value()
		}
	}
	return ""
}

func (this *rawConfigParser) GetListValue(section, option string) []string {
	this.RLock()
	defer this.RUnlock()

	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.ListValue()
		}
	}
	return nil
}