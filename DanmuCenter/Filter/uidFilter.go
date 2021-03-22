package Filter

import "github.com/Holynnchen/BiliBan2/DanmuCenter"

// uid小于该区域视作安全
type uidFilter struct {
	uidTarget int64
}

func (filter *uidFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.UserID < filter.uidTarget, ""
}
func (filter *uidFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}
func (filter *uidFilter) SafeCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}

// Filter-> uid<uidTarget
func NewUIDFilter(uidTarget int64) *uidFilter {
	return &uidFilter{
		uidTarget: uidTarget,
	}
}
