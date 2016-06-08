package config

type Option struct {
	name  string
	iv    string
	value []string
}

func NewOption(name, iv, value string) *Option {
	var opt = &Option{}
	opt.name = name
	opt.iv = iv
	opt.value = []string{value}
	return opt
}

func (this *Option) Name() string {
	return this.name
}

func (this *Option) Value() string {
	if len(this.value[0]) > 0 {
		return this.value[0]
	}
	return ""
}

func (this *Option) ListValue() []string {
	return this.value
}