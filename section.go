package config

type section struct {
	name    string
	options map[string]*Option
}
