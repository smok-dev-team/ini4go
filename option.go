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