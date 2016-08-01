package config

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/smartwalle/container"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
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
type Config struct {
	rawConfigParser
}

//func New() *Config {
//	return NewConfigWithBlock(true, true)
//}

func New(block bool) *Config {
	var c = &Config{}
	c.block = block
	c.mutex = &sync.RWMutex{}
	c.init()
	return c
}

////////////////////////////////////////////////////////////////////////////////
type rawConfigParser struct {
	mutex       *sync.RWMutex
	sectionKeys []string
	sections    map[string]*Section
	block       bool
}

func (this *rawConfigParser) Lock() {
	if this.block {
		this.mutex.Lock()
	}
}

func (this *rawConfigParser) Unlock() {
	if this.block {
		this.mutex.Unlock()
	}
}

func (this *rawConfigParser) RLock() {
	if this.block {
		this.mutex.RLock()
	}
}

func (this *rawConfigParser) RUnlock() {
	if this.block {
		this.mutex.RUnlock()
	}
}

func (this *rawConfigParser) init() {
	this.sectionKeys = nil
	this.sections = make(map[string]*Section)
}

func (this *rawConfigParser) Load(dir string) error {
	var fileInfo, err = os.Stat(dir)
	if err != nil {
		return err
	}

	var pathList []string

	if fileInfo.IsDir() {
		var file *os.File
		file, err = os.Open(dir)
		if err != nil {
			return err
		}

		var names []string
		names, err = file.Readdirnames(-1)

		file.Close()
		if err != nil {
			return err
		}

		for _, name := range names {
			var filePath = path.Join(dir, name)
			fileInfo, err = os.Stat(filePath)
			if err != nil {
				continue
			}

			if !fileInfo.IsDir() {
				pathList = append(pathList, filePath)
			}
		}
	} else {
		pathList = append(pathList, dir)
	}

	return this.LoadFiles(pathList...)
}

func (this *rawConfigParser) LoadFiles(files ...string) error {
	this.Lock()
	defer this.Unlock()

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
			currentSection = this.newSection(sectionName)
			currentSection.comments = append(currentSection.comments, comments...)
			comments = nil
			continue
		}

		if currentSection == nil {
			currentSection = this.newSection(strings.ToLower(k_DEFAULT_SECTION))
		}

		var optName, optIV, optValue = getOptionAndValue(sLine)
		optName = strings.TrimSpace(optName)
		optIV = strings.TrimSpace(optIV)
		optValue = strings.TrimSpace(optValue)

		if optName != "" {
			var opt = currentSection.newOption(optName, optIV)
			opt.AddValue(optValue)
			opt.AddComment(comments...)
			comments = nil
		}
	}
	return nil
}

func (this *rawConfigParser) WriteToFile(file string) error {
	var f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	err = this.writeTo(f)
	f.Close()
	return err
}

func (this *rawConfigParser) writeTo(w io.Writer) error {
	this.Lock()
	defer this.Unlock()

	var writer = bufio.NewWriter(w)

	for _, sectionName := range this.sectionKeys {

		var section = this.section(sectionName)
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
			for _, value := range opt.values {
				writer.WriteString(fmt.Sprintf("%s %s %s\n", opt.key, opt.iv, value))
			}
		}
		writer.WriteByte('\n')
	}
	return writer.Flush()
}

////////////////////////////////////////////////////////////////////////////////
func (this *rawConfigParser) Reset() {
	this.Lock()
	this.Unlock()
	this.init()
}

func (this *rawConfigParser) newSection(name string) *Section {
	var section = this.sections[name]
	if section == nil {
		section = NewSection(name)
		this.sections[name] = section
		this.sectionKeys = append(this.sectionKeys, name)
	}
	return section
}

func (this *rawConfigParser) NewSection(name string) *Section {
	this.Lock()
	defer this.Unlock()

	return this.newSection(name)
}

func (this *rawConfigParser) mustSection(name string) *Section {
	return this.newSection(name)
}

func (this *rawConfigParser) MustSection(name string) *Section {
	return this.NewSection(name)
}

func (this *rawConfigParser) section(name string) *Section {
	var s, _ = this.sections[name]
	return s
}

func (this *rawConfigParser) Section(name string) *Section {
	this.RLock()
	defer this.RUnlock()
	return this.section(name)
}

func (this *rawConfigParser) SectionNames() []string {
	this.RLock()
	defer this.RUnlock()

	var names = make([]string, len(this.sectionKeys))
	copy(names, this.sectionKeys)
	return names
}

func (this *rawConfigParser) SectionList() []*Section {
	this.RLock()
	defer this.RUnlock()

	var sList = make([]*Section, 0, len(this.sections))
	for _, value := range this.sections {
		sList = append(sList, value)
	}
	return sList
}

func (this *rawConfigParser) HasSection(section string) bool {
	this.RLock()
	defer this.RUnlock()

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
func (this *rawConfigParser) mustOption(section, option string) *Option {
	var s = this.mustSection(section)
	var opt = s.MustOption(option)
	return opt
}

func (this *rawConfigParser) MustOption(section, option string) *Option {
	this.Lock()
	defer this.Unlock()

	return this.mustOption(section, option)
}

func (this *rawConfigParser) option(section, option string) *Option {
	var s = this.section(section)
	if s != nil {
		return s.Option(option)
	}
	return nil
}

func (this *rawConfigParser) Option(section, option string) *Option {
	this.Lock()
	defer this.Unlock()
	return this.option(section, option)
}

func (this *rawConfigParser) Options(section string) []string {
	this.RLock()
	defer this.RUnlock()

	var s = this.section(section)
	if s != nil {
		return s.OptionKeys()
	}
	return nil
}

func (this *rawConfigParser) OptionList(section string) []*Option {
	this.RLock()
	defer this.RUnlock()

	var s = this.section(section)
	if s != nil {
		return s.OptionList()
	}
	return nil
}

func (this *rawConfigParser) HasOption(section, option string) bool {
	this.RLock()
	defer this.RUnlock()

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
func (this *rawConfigParser) SetValue(section, option, value string) {
	this.Lock()
	defer this.Unlock()

	var s = this.newSection(section)
	var opt = s.newOption(option, "=")
	opt.SetValue(value)
}

func (this *rawConfigParser) SetString(section, option, value string) {
	this.SetValue(section, option, value)
}

func (this *rawConfigParser) SetInt(section, option string, value int) {
	this.SetValue(section, option, fmt.Sprintf("%d", value))
}

func (this *rawConfigParser) SetInt64(section, option string, value int64) {
	this.SetValue(section, option, fmt.Sprintf("%d", value))
}

func (this *rawConfigParser) SetFloat32(section, option string, value float32) {
	this.SetValue(section, option, fmt.Sprintf("%f", value))
}

func (this *rawConfigParser) SetFloat64(section, option string, value float64) {
	this.SetValue(section, option, fmt.Sprintf("%f", value))
}

func (this *rawConfigParser) SetBool(section, option string, value bool) {
	this.SetValue(section, option, fmt.Sprintf("%t", value))
}

////////////////////////////////////////////////////////////////////////////////
func (this *rawConfigParser) GetValue(section, option string) string {
	return this.MustValue(section, option, "")
}

func (this *rawConfigParser) MustValue(section, option, defaultValue string) string {
	this.RLock()
	defer this.RUnlock()
	var opt = this.mustOption(section, option)
	return opt.MustString(defaultValue)
}

func (this *rawConfigParser) MustInt(section, option string, defaultValue int) int {
	this.RLock()
	defer this.RUnlock()

	var opt = this.mustOption(section, option)
	return opt.MustInt(defaultValue)
}

func (this *rawConfigParser) MustInt64(section, option string, defaultValue int64) int64 {
	this.RLock()
	defer this.RUnlock()

	var opt = this.mustOption(section, option)
	return opt.MustInt64(defaultValue)
}

func (this *rawConfigParser) MustFloat32(section, option string, defaultValue float32) float32 {
	this.RLock()
	defer this.RUnlock()

	var opt = this.mustOption(section, option)
	return opt.MustFloat32(defaultValue)
}

func (this *rawConfigParser) MustFloat64(section, option string, defaultValue float64) float64 {
	this.RLock()
	defer this.RUnlock()

	var opt = this.mustOption(section, option)
	return opt.MustFloat64(defaultValue)
}

func (this *rawConfigParser) MustBool(section, option string, defaultValue bool) bool {
	this.RLock()
	defer this.RUnlock()

	var opt = this.mustOption(section, option)
	return opt.MustBool(defaultValue)
}

func (this *rawConfigParser) GetValues(section, option string) []string {
	this.RLock()
	defer this.RUnlock()

	var s = this.sections[section]
	if s != nil {
		var opt = s.options[option]
		if opt != nil {
			return opt.Values()
		}
	}
	return nil
}