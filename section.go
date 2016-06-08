package config

import "github.com/smartwalle/going/container"

type Section struct {
	name       string
	optionKeys []string
	options    map[string]*Option
}

func NewSection(name string) *Section {
	var section = &Section{}
	section.name = name
	section.options = make(map[string]*Option)
	return section
}

func (this *Section) Name() string {
	return this.name
}

func (this *Section) NewOption(key, iv, value string) {
	var opt = this.options[key]
	if opt == nil {
		opt = NewOption(key, iv, value)
		this.options[key] = opt
		this.optionKeys = append(this.optionKeys, key)
	} else {
		opt.value = append(opt.value, value)
	}
}

func (this *Section) RemoveOption(key string) {
	delete(this.options, key)
	container.Remove(&this.optionKeys, key)
}

func (this *Section) HasOption(key string) bool {
	var _, ok = this.options[key]
	return ok
}

func (this *Section) Option(key string) *Option {
	var opt = this.options[key]
	return opt
}

func (this *Section) OptionKeys() []string {
	return this.optionKeys
}

func (this *Section) OptionList() []*Option {
	var oList = make([]*Option, 0, len(this.options))
	for _, value := range this.options {
		oList = append(oList, value)
	}
	return oList
}