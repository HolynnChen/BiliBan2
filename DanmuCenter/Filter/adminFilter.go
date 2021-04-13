package Filter

import (
	"github.com/Holynnchen/BiliBan2/DanmuCenter"
)

// 忽略房管
type AdminFilter struct{}

func (*AdminFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return danmu.IsAdmin, ""
}

func NewAdminFilter() *AdminFilter {
	return &AdminFilter{}
}
