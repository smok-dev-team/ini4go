package config

type Option struct {
	key   string
	iv    string
	value []string
}

func NewOption(key, iv, value string) *Option {
	var opt = &Option{}
	opt.key = key
	opt.iv = iv
	opt.value = []string{value}
	return opt
}

func (this *Option) Key() string {
	return this.key
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