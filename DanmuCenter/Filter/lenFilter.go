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
	levelMap            map[int]int // 特定等级特定处理
}

type lenFilterOption func(*LenFilter)

func (filter *LenFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	content := danmu.Content
	lenTarget := filter.lenTarget
	if filter.repeatGroupCompress != nil {
		content = filter.repeatGroupCompress(content)
	}
	if specialLentarget, ok := filter.levelMap[danmu.UserLevel]; ok {
		lenTarget = specialLentarget
	}
	return utf8.RuneCountInString(content) < lenTarget, ""
}

func NewLenFilter(lenTarget int, options ...lenFilterOption) *LenFilter {
	filter := &LenFilter{
		lenTarget: lenTarget,
		levelMap:  make(map[int]int),
	}
	for _, option := range options {
		option(filter)
	}
	return filter
}

func SetLenFilterCompressRepeatGroup(minLen int) lenFilterOption {
	return func(lf *LenFilter) {
		lf.repeatGroupCompress = Utils.CompressRepeatGroup2(minLen)
	}
}

func SetSpecialLevelLenTarget(level, lenTarget int) lenFilterOption {
	return func(lf *LenFilter) {
		lf.levelMap[level] = lenTarget
	}
}
