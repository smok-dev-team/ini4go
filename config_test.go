package config

import (
	"testing"
	"fmt"
)

func TestSectionNameRegex(t *testing.T) {
	if getSectionName("[section-1]") != "section-1" {
		t.Error("不能正常解析 section-1")
	}

	if getSectionName("section-1") == "section-1" {
		t.Error("解析 section 异常")
	}

	if getSectionName("[支持中文]") != "支持中文" {
		t.Error("不能正常解析中文")
	}
}

func TestOptionIsExist(t *testing.T) {
	var r = NewConfig()
	r.LoadFiles("./test.conf")
	fmt.Println(r.HasOption("default", "sk2"))
}

func TestLoadFile(t *testing.T) {
	var r = NewConfig()
	r.LoadFiles("./test.conf")

	fmt.Println(r.GetValue("default", "sk2"))
	fmt.Println(r.GetValue("s2","sk22"))
}

func TestOutput(t *testing.T) {
	var r = NewConfig()
	r.LoadFiles("./PerfStringBackup.ini", "./test.conf")


	fmt.Println(r.SectionNames())
	fmt.Println(r.Options("s2"))

	fmt.Println(r.WriteToFile("./a.conf"))
}

func TestNew(t *testing.T) {
	var r = NewConfig()
	r.SetValue("s1", "k1", "v1")
	r.SetValue("s1", "k2", "v2")
	r.SetValue("s2", "k1", "v1")

	r.MustOption("s3", "kk").AppendValue("sdfsf", "Ser", "xc")

	r.WriteToFile("./a.ini")
}