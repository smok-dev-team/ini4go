package config

import (
	"testing"
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
//
//func TestOutput(t *testing.T) {
//	var r = &RawConfigParser{}
//	r.LoadFiles("./PerfStringBackup.ini", "./test.conf")
//
//	//var sList = r.Sections()
//
//	fmt.Println(r.SectionNames())
//	fmt.Println(r.Options("s2"))
//
//	fmt.Println(r.WriteToFile("./a.conf"))
//}

func TestNew(t *testing.T) {
	var r = NewConfig()
	r.SetValue("s1", "k1", "v1")
	r.SetValue("s1", "k2", "v2")
	r.SetValue("s2", "k1", "v1")

	r.WriteToFile("./a.ini")



//	r.LoadFiles("./a.ini")
//
//	fmt.Println(r.Section("s1").Comment())
//	fmt.Println(r.Option("s1", "k2").Comment())

}