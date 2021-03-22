package Filter

import "github.com/Holynnchen/BiliBan2/DanmuCenter"

// userLevel >= levelTarget视为正常
type levelFilter struct {
	levelTarget int
}

func (filter *levelFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.UserLevel >= filter.levelTarget, ""
}
func (filter *levelFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}
func (filter *levelFilter) SafeCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}

// Filter-> userLever >= levelTarget
func NewUserLevelFilter(levelTarget int) *levelFilter {
	return &levelFilter{
		levelTarget: levelTarget,
	}
}
