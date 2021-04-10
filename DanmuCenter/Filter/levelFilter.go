package Filter

import "github.com/Holynnchen/BiliBan2/DanmuCenter"

// userLevel >= levelTarget视为正常
type LevelFilter struct {
	levelTarget int
}

func (filter *LevelFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.UserLevel >= filter.levelTarget, ""
}

// Filter-> userLever >= levelTarget
func NewUserLevelFilter(levelTarget int) *LevelFilter {
	return &LevelFilter{
		levelTarget: levelTarget,
	}
}
