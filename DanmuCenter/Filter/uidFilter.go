package Filter

import "github.com/Holynnchen/BiliBan2/DanmuCenter"

// uid小于该区域视作安全
type UidFilter struct {
	uidTarget int64
}

func (filter *UidFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.UserID < filter.uidTarget, ""
}

// Filter-> uid<uidTarget
func NewUIDFilter(uidTarget int64) *UidFilter {
	return &UidFilter{
		uidTarget: uidTarget,
	}
}
