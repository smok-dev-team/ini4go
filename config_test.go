package config

import (
	"testing"
	"fmt"
)

//func TestSectionNameRegex(t *testing.T) {
//	if getSectionName("[section-1]") != "section-1" {
//		t.Error("不能正常解析 section-1")
//	}
//
//	if getSectionName("section-1") == "section-1" {
//		t.Error("解析 section 异常")
//	}
//
//	if getSectionName("[支持中文]") != "支持中文" {
//		t.Error("不能正常解析中文")
//	}
//}
//
//func TestOptionIsExist(t *testing.T) {
//	var r = &RawConfigParser{}
//	r.LoadFile("./test.conf")
//	fmt.Println(r.HasOption("default", "sk2"))
//}
//
//func TestLoadFile(t *testing.T) {
//	var r = &RawConfigParser{}
//	r.LoadFile("./test.conf")
//
//	fmt.Println(r.GetListValue("default", "sk2"))
//}

func TestOutput(t *testing.T) {
	var r = &RawConfigParser{}
	r.LoadFiles("./PerfStringBackup.ini", "./test.conf")

	//var sList = r.Sections()

	fmt.Println(r.SectionNames())
	fmt.Println(r.Options("s2"))

	//for _, name := range sList {
	//	var s = r.Section(name)
	//	var oList = s.Options()
	//	for _, oName := range oList {
	//		var opt = s.Option(oName)
	//		fmt.Println(opt.Name(), opt.Value())
	//	}
	//}
}