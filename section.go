package config

type Section struct {
	name    string
	options map[string]*Option
}

func NewSection(name string) *Section {
	var section = &Section{}
	section.name = name
	section.options = make(map[string]*Option)
	return section
}

func (this *Section) NewOption(name, iv, value string) {
	var opt = this.options[name]
	if opt == nil {
		opt = NewOption(name, iv, value)
		this.options[name] = opt
	} else {
		opt.value = append(opt.value, value)
	}
}

func (this *Section) RemoveOption(name string) {
	delete(this.options, name)
}

func (this *Section) HasOption(name string) bool {
	var _, ok = this.options[name]
	return ok
}

func (this *Section) Option(name string) *Option {
	var opt = this.options[name]
	return opt
}

func (this *Section) Options() []string {
	var oList = make([]string, 0, len(this.options))
	for key := range this.options {
		oList = append(oList, key)
	}
	return oList
}