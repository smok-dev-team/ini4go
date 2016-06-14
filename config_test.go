package config

import (
	"testing"
	"fmt"
	"time"
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

	if r.HasOption("s1", "sk2") == false {
		t.Error("s1 中有 sk2")
	}
}

func TestLoadFile(t *testing.T) {
	var r = NewConfig()
	r.LoadFiles("./test.conf")

	if r.GetValue("default", "dk1") != "dkv1" {
		t.Error("default -> dk1 应该为 dkv1")
	}

	if r.GetValue("s1", "sk1") != "skv1" {
		t.Error("s1 -> sk1 应该为 skv1")
	}

	if r.GetValue("不存在的section", "不存在的option") != "" {
		t.Error("不存在的section")
	}

	if r.MustValue("不存在的section", "不存在的option", "但是有默认值") != "但是有默认值" {
		t.Error("不存在的section, 但是有默认值的哦")
	}
}

func TestOutput(t *testing.T) {
	var r = NewConfig()
	r.SetValue("s1", "p1", "v1")
	r.MustSection("s1").MustOption("p2").SetValue("v2")
	r.MustSection("s2").MustOption("p2").SetValue("v2")
	fmt.Println(r.WriteToFile("./output.conf"))
}

func TestAppend(t *testing.T) {
	var r = NewConfig()
	r.SetValue("s1", "k1", "v1")
	r.SetValue("s1", "k2", "v2")
	r.SetValue("s2", "k1", "v1")

	r.MustOption("s3", "k1").AddValue("第一个值")
	r.MustOption("s3", "k1").AddValue("第二个值", "第三个值", "第四个值")

	fmt.Println(r.MustValue("s3", "k1", "oh no"))
	fmt.Println(r.MustOption("s3", "k1").Values())

	r.MustOption("s4", "k1").SetBool(true)
	r.MustOption("s4", "k2").SetInt(11)
	r.MustOption("s4", "k3").SetInt64(61)
	r.MustOption("s4", "k4").SetFloat32(111.111)
	r.MustOption("s4", "k5").SetTime(time.Now())
	r.MustOption("s4", "k6").SetTimeWithLayout(time.Now(), "2006-01-02 15:04:05")


	fmt.Println(r.MustOption("s4", "k6").TimeWithLayout("2006-01-02 15:04:05"))

	r.WriteToFile("./output2.conf")
}

//func TestLoad(t *testing.T) {
//	var r = NewConfig()
//	r.LoadFiles("./PerfStringBackup.ini")
//
//	var sectionNames = r.SectionNames()
//	for _, name := range sectionNames {
//		fmt.Println(name)
//		var section = r.Section(name)
//		var optKeys = section.OptionKeys()
//		for _, key := range optKeys {
//			var opt = section.Option(key)
//			fmt.Println("   ", opt.Key(), " = ", opt.Value())
//		}
//	}
//}