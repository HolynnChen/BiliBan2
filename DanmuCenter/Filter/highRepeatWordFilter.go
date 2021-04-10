package Filter

import (
	"unicode/utf8"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
)

//单重复率大于repeatTarget视作正常弹幕
type HighRepeatWordFilter struct {
	repeatTarget float32
}

func (filter *HighRepeatWordFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	countMap := make(map[rune]int)
	max := 0
	for _, data := range danmu.Content {
		countMap[data]++
		if countMap[data] > max {
			max = countMap[data]
		}
	}
	return float32(max)/float32(utf8.RuneCountInString(danmu.Content)) > filter.repeatTarget, ""
}

func NewHighReatWordFilter(repeatTarget float32) *HighRepeatWordFilter {
	return &HighRepeatWordFilter{
		repeatTarget: repeatTarget,
	}
}
