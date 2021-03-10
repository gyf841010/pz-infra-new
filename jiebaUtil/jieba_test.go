package jiebaUtil

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJieBaCutForSearch(t *testing.T) {
	Convey("Test JieBa Cut For Search", t, func() {
		Convey("Test Valid Email numbers", func() {
			fmt.Print("【搜索引擎模式】：")
			result := CutWordForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造")
			So(len(result), ShouldEqual, 15)
		})
	})
}
