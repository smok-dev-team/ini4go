package config

import (
	"testing"
)

func TestSectionNameRegex(t *testing.T) {
	if sectionName("[section-1]") != "section-1" {
		t.Error("不能正常解析 section-1")
	}

	if sectionName("section-1") == "section-1" {
		t.Error("解析 section 异常")
	}

	if sectionName("[支持中文]") != "支持中文" {
		t.Error("不能正常解析中文")
	}
}