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

func (this *Section) HasOption(name string) bool {
	var _, ok = this.options[name]
	return ok
}

func (this *Section) GetOption(name string) *Option {
	var opt = this.options[name]
	return opt
}