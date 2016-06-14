package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"regexp"
)

////////////////////////////////////////////////////////////////////////////////
var varRegexp = regexp.MustCompile(`%\((\S[^\(\)]+)\)s`)

func getVarName(src string) [][]string {
	var rList = varRegexp.FindAllStringSubmatch(src, -1)
	return rList
}

////////////////////////////////////////////////////////////////////////////////
type Option struct {
	section  *Section
	key      string
	iv       string
	values   []string
	comments []string
}

func NewOption(section *Section, key, iv string, values []string) *Option {
	var opt = &Option{}
	opt.key = key
	opt.iv = iv
	opt.values = values
	opt.section = section
	return opt
}

func (this *Option) Key() string {
	return this.key
}

func (this *Option) Comments() []string {
	return this.comments
}

func (this *Option) Comment() string {
	if len(this.comments) > 0 {
		return this.comments[0]
	}
	return ""
}

func (this *Option) AddComment(comment ...string) {
	if len(comment) > 0 {
		this.comments = append(this.comments, comment...)
	}
}

func (this *Option) parseValue(raw string) (result string) {
	//var varNameList = getVarName(raw)
	//result = raw
	//for _, l := range varNameList {
	//	var keyVar = l[0]
	//	var keyName = l[1]
	//	var value = this.section.MustOption(keyName).Value()
	//	result = strings.Replace(raw, keyVar, value, -1)
	//}
	//return result

	result = varRegexp.ReplaceAllStringFunc(raw, func (src string) string {
		var key = src[2:len(src)-2]
		var value = this.section.MustOption(key).Value()
		return value
	})
	return result
}

func (this *Option) Value() string {
	return this.ValueAt(0)
}

func (this *Option) ValueAt(index int) string {
	if len(this.values) > index {
		return this.parseValue(this.values[index])
	}
	return ""
}

func (this *Option) Values() []string {
	var l = len(this.values)
	var newValues = make([]string, l)

	for i:=0; i<l; i++ {
		newValues[i] = this.ValueAt(i)
	}

	return newValues
}

func (this *Option) SetValue(v string) {
	this.values = []string{v}
}

func (this *Option) AddValue(v ...string) {
	if len(v) > 0 {
		this.values = append(this.values, v...)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (this *Option) String() string {
	return this.Value()
}

func (this *Option) MustString(defaultValue string) string {
	if len(this.String()) > 0 {
		return this.String()
	}
	return defaultValue
}

func (this *Option) SetString(s string) {
	this.SetValue(s)
}

func (this *Option) Int() (int, error) {
	var v = this.String()
	return strconv.Atoi(v)
}

func (this *Option) MustInt(defaultValue int) int {
	var v, err = this.Int()
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetInt(v int) {
	var s = fmt.Sprintf("%d", v)
	this.SetValue(s)
}

func (this *Option) Int64() (int64, error) {
	var v = this.String()
	return strconv.ParseInt(v, 10, 64)
}

func (this *Option) MustInt64(defaultValue int64) int64 {
	var v, err = this.Int64()
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetInt64(v int64) {
	var s = fmt.Sprintf("%d", v)
	this.SetValue(s)
}

func (this *Option) Float32() (float32, error) {
	var v = this.String()
	var fv, err = strconv.ParseFloat(v, 32)
	return float32(fv), err
}

func (this *Option) MustFloat32(defaultValue float32) float32 {
	var v, err = this.Float32()
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetFloat32(v float32) {
	var s = fmt.Sprintf("%f", v)
	this.SetValue(s)
}

func (this *Option) Float64() (float64, error) {
	var v = this.String()
	return strconv.ParseFloat(v, 64)
}

func (this *Option) MustFloat64(defaultValue float64) float64 {
	var v, err = this.Float64()
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetFloat64(v float64) {
	var s = fmt.Sprintf("%f", v)
	this.SetValue(s)
}

func (this *Option) Bool() (bool, error) {
	var v = strings.ToLower(this.String())
	switch v {
	case "1", "true", "yes", "on", "t", "y":
		return true, nil
	case "0", "false", "no", "off", "f", "n":
		return false, nil
	}
	return false, errors.New(fmt.Sprintf("parsing \"%s\": invalid syntax", v))
}

func (this *Option) MustBool(defaultValue bool) bool {
	var v, err = this.Bool()
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetBool(v bool) {
	var s = fmt.Sprintf("%t", v)
	this.SetValue(s)
}

func (this *Option) Time() (time.Time, error) {
	return this.TimeWithLayout("2006-01-02 15:04:05.999999999 -0700 MST")
}

func (this *Option) MustTime(defaultValue time.Time) time.Time {
	var v, err = this.Time()
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetTime(v time.Time) {
	var s = v.String()
	this.SetValue(s)
}

func (this *Option) TimeWithLayout(layout string) (time.Time, error) {
	var v = this.String()
	return time.Parse(layout, v)
}

func (this *Option) MustTimeWithLayout(layout string, defaultValue time.Time) time.Time {
	var v, err = this.TimeWithLayout(layout)
	if err == nil {
		return v
	}
	return defaultValue
}

func (this *Option) SetTimeWithLayout(v time.Time, layout string) {
	var s = v.Format(layout)
	this.SetValue(s)
}
