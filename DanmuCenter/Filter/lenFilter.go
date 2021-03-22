package Filter

import (
	"unicode/utf8"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
)

//小于lenTarget视作正常弹幕
type lenFilter struct {
	repeatGroupCompress func(string) string
	lenTarget           int
}

type lenFilterOption func(*lenFilter)

func (filter *lenFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	content := danmu.Content
	if filter.repeatGroupCompress != nil {
		content = filter.repeatGroupCompress(content)
	}
	return utf8.RuneCountInString(content) < filter.lenTarget, ""
}

func NewLenFilter(lenTarget int, options ...lenFilterOption) *lenFilter {
	filter := &lenFilter{
		lenTarget: lenTarget,
	}
	for _, option := range options {
		option(filter)
	}
	return filter
}

func SetLenFilterCompressRepeatGroup(minLen int) lenFilterOption {
	return func(lf *lenFilter) {
		lf.repeatGroupCompress = Utils.CompressRepeatGroup(minLen)
	}
}
