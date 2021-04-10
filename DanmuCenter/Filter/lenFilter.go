package Filter

import (
	"unicode/utf8"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
)

//小于lenTarget视作正常弹幕
type LenFilter struct {
	repeatGroupCompress func(string) string
	lenTarget           int
}

type lenFilterOption func(*LenFilter)

func (filter *LenFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	content := danmu.Content
	if filter.repeatGroupCompress != nil {
		content = filter.repeatGroupCompress(content)
	}
	return utf8.RuneCountInString(content) < filter.lenTarget, ""
}

func NewLenFilter(lenTarget int, options ...lenFilterOption) *LenFilter {
	filter := &LenFilter{
		lenTarget: lenTarget,
	}
	for _, option := range options {
		option(filter)
	}
	return filter
}

func SetLenFilterCompressRepeatGroup(minLen int) lenFilterOption {
	return func(lf *LenFilter) {
		lf.repeatGroupCompress = Utils.CompressRepeatGroup(minLen)
	}
}
