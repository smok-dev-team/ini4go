package ini4go

import "sync"

type Section struct {
	name       string
	optionKeys []string
	options    sync.Map
	comments   []string
}

func NewSection(name string) *Section {
	var section = &Section{}
	section.name = name
	section.options = sync.Map{}
	return section
}

func (this *Section) Name() string {
	return this.name
}

func (this *Section) Comments() []string {
	return this.comments
}

func (this *Section) Comment() string {
	if len(this.comments) > 0 {
		return this.comments[0]
	}
	return ""
}

func (this *Section) AddComment(comment string) {
	this.comments = append(this.comments, comment)
}

func (this *Section) newOption(key, iv string) *Option {
	opt, _ := this.options.Load(key)
	if opt == nil {
		opt = NewOption(this, key, iv, nil)
		this.options.LoadOrStore(key, opt)
		this.optionKeys = append(this.optionKeys, key)
	}
	return opt.(*Option)
}

func (this *Section) NewOption(key, iv string, value, comments []string) *Option {
	opt, _ := this.options.Load(key)
	if opt == nil {
		opt = NewOption(this, key, iv, nil)
		this.options.LoadOrStore(key, opt)
		this.optionKeys = append(this.optionKeys, key)
	}

	opt.(*Option).AddValue(value...)
	opt.(*Option).AddComment(comments...)
	return opt.(*Option)
}

func (this *Section) RemoveOption(key string) {
	this.options.Delete(key)
	var index = -1
	for i, opt := range this.optionKeys {
		if opt == key {
			index = i
			break
		}
	}

	if index >= 0 {
		this.optionKeys = append(this.optionKeys[0:index], this.optionKeys[index+1:]...)
	}
}

func (this *Section) HasOption(key string) bool {
	_, ok := this.options.Load(key)
	return ok
}

func (this *Section) MustOption(key string) *Option {
	var opt = this.newOption(key, "=")
	return opt
}

func (this *Section) Option(key string) *Option {
	opt, _ := this.options.Load(key)
	return opt.(*Option)
}

func (this *Section) OptionKeys() []string {
	var keys = make([]string, len(this.optionKeys))
	copy(keys, this.optionKeys)
	return keys
}

func (this *Section) OptionList() []*Option {

	var oList = make([]*Option, 0)
	// 遍历 this.options
	f := func(key, value interface{}) bool {
		oList = append(oList, value.(*Option))
		return true
	}
	this.options.Range(f)
	return oList
}
