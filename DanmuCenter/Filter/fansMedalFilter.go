package Filter

import (
	"github.com/Holynnchen/BiliBan2/DanmuCenter"
)

// fansMedalLever >=levelTarget视为正常
type FansMedalFilter struct {
	levelTarget int
}

func (filter *FansMedalFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.MedalLevel >= filter.levelTarget, ""
}

// Filter-> fansMedalLevel >= leverTarget
func NewFansMedalFilter(levelTarget int) *FansMedalFilter {
	return &FansMedalFilter{
		levelTarget: levelTarget,
	}
}
