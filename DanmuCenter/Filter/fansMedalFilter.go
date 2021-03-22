package Filter

import (
	"github.com/Holynnchen/BiliBan2/DanmuCenter"
)

// fansMedalLever >=levelTarget视为正常
type fansMedalFilter struct {
	levelTarget int
}

func (filter *fansMedalFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.MedalLevel >= filter.levelTarget, ""
}
func (filter *fansMedalFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}
func (filter *fansMedalFilter) SafeCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}

// Filter-> fansMedalLevel >= leverTarget
func NewFansMedalFilter(levelTarget int) *fansMedalFilter {
	return &fansMedalFilter{
		levelTarget: levelTarget,
	}
}
